package istio

import (
	"context"
	"fmt"
	"sort"

	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	gloov1 "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"
	"github.com/solo-io/supergloo/pkg/api/external/istio/networking/v1alpha3"
	"github.com/solo-io/supergloo/pkg/api/v1"
)

type RoutingSyncer struct {
	WriteSelector             map[string]string // for reconciling only our resources
	WriteNamespace            string
	DestinationRuleReconciler v1alpha3.DestinationRuleReconciler
	VirtualServiceReconciler  v1alpha3.VirtualServiceReconciler
	Reporter                  reporter.Reporter
}

func (s *RoutingSyncer) Sync(ctx context.Context, snap *v1.TranslatorSnapshot) error {
	ctx = contextutils.WithLogger(ctx, "routing-syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v meshes, %v upstreams)", snap.Hash(),
		len(snap.Meshes), len(snap.Upstreams))
	defer logger.Infof("end sync %v", snap.Hash())
	logger.Debugf("%v", snap)

	destinationRules := createDestinationRules(false, snap.Upstreams.List())
	virtualServices, err := createVirtualServices(snap.Meshes.List(), snap.Upstreams.List())
	if err != nil {
		return errors.Wrapf(err, "creating virtual services from snapshot")
	}
	for _, res := range destinationRules {
		resources.UpdateMetadata(res, func(meta *core.Metadata) {
			meta.Namespace = s.WriteNamespace
			if meta.Annotations == nil {
				meta.Annotations = make(map[string]string)
			}
			meta.Annotations["created_by"] = "supergloo"
			for k, v := range s.WriteSelector {
				meta.Labels[k] = v
			}
		})
	}
	for _, res := range virtualServices {
		resources.UpdateMetadata(res, func(meta *core.Metadata) {
			meta.Namespace = s.WriteNamespace
			if meta.Annotations == nil {
				meta.Annotations = make(map[string]string)
			}
			if meta.Labels == nil {
				meta.Labels = make(map[string]string)
			}
			meta.Annotations["created_by"] = "supergloo"
			for k, v := range s.WriteSelector {
				meta.Labels[k] = v
			}
		})
	}
	return s.writeIstioCrds(ctx, destinationRules, virtualServices)
}

func (s *RoutingSyncer) writeIstioCrds(ctx context.Context, destinationRules v1alpha3.DestinationRuleList, virtualServices v1alpha3.VirtualServiceList) error {
	opts := clients.ListOpts{
		Ctx:      ctx,
		Selector: s.WriteSelector,
	}
	contextutils.LoggerFrom(ctx).Infof("reconciling %v destination rules", len(destinationRules))
	if err := s.DestinationRuleReconciler.Reconcile(s.WriteNamespace, destinationRules, preserveDestinationRule, opts); err != nil {
		return errors.Wrapf(err, "reconciling destination rules")
	}
	contextutils.LoggerFrom(ctx).Infof("reconciling %v virtual services", len(virtualServices))
	if err := s.VirtualServiceReconciler.Reconcile(s.WriteNamespace, virtualServices, preserveVirtualService, opts); err != nil {
		return errors.Wrapf(err, "reconciling virtual services")
	}
	return nil
}

func preserveDestinationRule(original, desired *v1alpha3.DestinationRule) (bool, error) {
	original.Metadata = desired.Metadata
	original.Status = desired.Status
	return !proto.Equal(original, desired), nil
}

func preserveVirtualService(original, desired *v1alpha3.VirtualService) (bool, error) {
	original.Metadata = desired.Metadata
	original.Status = desired.Status
	return !proto.Equal(original, desired), nil
}

func subsetName(us *gloov1.Upstream) string {
	return fmt.Sprintf("%v.%v", us.Metadata.Namespace, us.Metadata.Name)
}

func createDestinationRules(enableTls bool, upstreams gloov1.UpstreamList) v1alpha3.DestinationRuleList {
	subsetsByDestination := make(map[string][]*v1alpha3.Subset)
	// only support kube upstreams for now
	for _, us := range upstreams {
		switch specType := us.UpstreamSpec.UpstreamType.(type) {
		case *gloov1.UpstreamSpec_Kube:
			if len(specType.Kube.Selector) == 0 {
				// no need for a subset
				continue
			}
			host := fmt.Sprintf("%v.%v.svc.cluster.local", specType.Kube.ServiceName, specType.Kube.ServiceNamespace)
			subsetsByDestination[host] = append(subsetsByDestination[host], &v1alpha3.Subset{
				Name:   subsetName(us),
				Labels: specType.Kube.Selector,
			})
		}
	}

	// TODO ilackarms: make enableTls a per-mesh variable
	var trafficPolicy *v1alpha3.TrafficPolicy
	if enableTls {
		trafficPolicy = &v1alpha3.TrafficPolicy{
			Tls: &v1alpha3.TLSSettings{
				Mode: v1alpha3.TLSSettings_ISTIO_MUTUAL,
			},
		}
	}
	var destinationRules v1alpha3.DestinationRuleList
	for host, subsets := range subsetsByDestination {
		destinationRules = append(destinationRules, &v1alpha3.DestinationRule{
			Metadata: core.Metadata{
				Name: host,
			},
			Host:          host,
			Subsets:       subsets,
			TrafficPolicy: trafficPolicy,
		})
	}
	sort.SliceStable(destinationRules, func(i, j int) bool {
		return destinationRules[i].Host < destinationRules[j].Host
	})
	return destinationRules
}

func getHostsForUpstream(us *gloov1.Upstream) ([]string, error) {
	switch specType := us.UpstreamSpec.UpstreamType.(type) {
	case *gloov1.UpstreamSpec_Aws:
		return nil, errors.Errorf("aws not implemented")
	case *gloov1.UpstreamSpec_Azure:
		return nil, errors.Errorf("azure not implemented")
	case *gloov1.UpstreamSpec_Kube:
		return []string{
			fmt.Sprintf("%v.%v.svc.cluster.local", specType.Kube.ServiceName, specType.Kube.ServiceNamespace),
			specType.Kube.ServiceName,
		}, nil
	case *gloov1.UpstreamSpec_Static:
		var hosts []string
		for _, h := range specType.Static.Hosts {
			hosts = append(hosts, h.Addr)
		}
		return hosts, nil
	}
	return nil, errors.Errorf("unsupported upstream type %v", us)
}

// TODO ilackarms: currently unused
func getPortForUpstream(us *gloov1.Upstream) (uint32, error) {
	switch specType := us.UpstreamSpec.UpstreamType.(type) {
	case *gloov1.UpstreamSpec_Aws:
		return 0, errors.Errorf("aws not implemented")
	case *gloov1.UpstreamSpec_Azure:
		return 0, errors.Errorf("azure not implemented")
	case *gloov1.UpstreamSpec_Kube:
		return specType.Kube.ServicePort, nil
	case *gloov1.UpstreamSpec_Static:
		// TODO(ilackarms): handle cases where port changes between hosts
		for _, h := range specType.Static.Hosts {
			return h.Port, nil
		}
		return 0, errors.Errorf("no hosts found on static upstream")
	}
	return 0, errors.Errorf("unknown upstream type")
}

func createVirtualServices(meshes v1.MeshList, upstreams gloov1.UpstreamList) (v1alpha3.VirtualServiceList, error) {
	var virtualServices v1alpha3.VirtualServiceList
	for _, mesh := range meshes {
		if mesh.TargetMesh == nil {
			return nil, errors.Errorf("invalid mesh %v: target_mesh required", mesh.Metadata.Ref())
		}
		if mesh.TargetMesh.MeshType != v1.MeshType_ISTIO {
			continue
		}
		if mesh.Routing == nil {
			continue
		}

		ingressRoutesByDomain := make(map[string][]*v1.Route)
		meshRoutesByDomain := make(map[string][]*v1.Route)
		for _, route := range mesh.Routing.Routes {
			if len(route.Domains) == 0 {
				route.Domains = []string{"*"}
			}
			for _, domain := range route.Domains {
				if route.EnabledForIngress {
					ingressRoutesByDomain[domain] = append(ingressRoutesByDomain[domain], route)
				}
				if route.EnabledForIngress {
					meshRoutesByDomain[domain] = append(meshRoutesByDomain[domain], route)
				}
			}
		}

		for domain, routes := range ingressRoutesByDomain {
			var istioRoutes []*v1alpha3.HTTPRoute
			for _, route := range routes {
				istioRoute, err := convertRoute(route, upstreams)
				if err != nil {
					return nil, errors.Wrapf(err, "invalid route %v", route)
				}
				istioRoutes = append(istioRoutes, istioRoute)
			}
			vs := &v1alpha3.VirtualService{
				Metadata: core.Metadata{
					Name: "supergloo-" + domain,
				},
				// in istio api, this is equivalent to []string{"mesh"}
				// which includes all pods in the mesh, with no selectors
				// and no ingresses
				Gateways: []string{"gateway"},
				Hosts:    []string{domain},
				Http:     istioRoutes,
			}
			virtualServices = append(virtualServices, vs)
		}
	}
	sort.SliceStable(virtualServices, func(i, j int) bool {
		return virtualServices[i].Metadata.Less(virtualServices[j].Metadata)
	})
	return virtualServices, nil
}

func convertRoute(route *v1.Route, upstreams gloov1.UpstreamList) (*v1alpha3.HTTPRoute, error) {
	istioRoute, err := convertAction(route.Action, upstreams)
	if err != nil {
		return nil, errors.Wrapf(err, "converting route action %v", route.Action)
	}
	var mirror *v1alpha3.Destination
	if route.Mirror != nil {
		mirror, err = convertDestination(route.Mirror, upstreams)
		if err != nil {
			return nil, errors.Wrapf(err, "converting route mirror %v", route.Mirror)
		}
	}
	return &v1alpha3.HTTPRoute{
		Match:                 convertMatch(route.RequestMatchers),
		Route:                 istioRoute,
		Timeout:               route.Timeout,
		Retries:               convertRetry(route.Retries),
		Fault:                 convertFault(route.Fault),
		Mirror:                mirror,
		CorsPolicy:            convertCorsPolicy(route.CorsPolicy),
		RemoveResponseHeaders: route.RemoveResponseHeaders,
		AppendResponseHeaders: route.AppendResponseHeaders,
		RemoveRequestHeaders:  route.RemoveRequestHeaders,
	}, nil
}

func convertRetry(retry *v1.HTTPRetry) *v1alpha3.HTTPRetry {
	if retry == nil {
		return nil
	}
	return &v1alpha3.HTTPRetry{
		Attempts:      retry.Attempts,
		PerTryTimeout: retry.PerTryTimeout,
	}
}
func convertFault(fault *v1.HTTPFaultInjection) *v1alpha3.HTTPFaultInjection {
	if fault == nil {
		return nil
	}
	var delay *v1alpha3.HTTPFaultInjection_Delay
	if fault.Delay != nil {
		delay = &v1alpha3.HTTPFaultInjection_Delay{
			Percentage: convertPercentage(fault.Delay.Percentage),
		}
		if fault.Delay.HttpDelayType != nil {
			switch delayType := fault.Delay.HttpDelayType.(type) {
			case *v1.HTTPFaultInjection_Delay_FixedDelay:
				delay.HttpDelayType = &v1alpha3.HTTPFaultInjection_Delay_FixedDelay{
					FixedDelay: delayType.FixedDelay,
				}
			case *v1.HTTPFaultInjection_Delay_ExponentialDelay:
				delay.HttpDelayType = &v1alpha3.HTTPFaultInjection_Delay_ExponentialDelay{
					ExponentialDelay: delayType.ExponentialDelay,
				}
			}
		}
	}
	var abort *v1alpha3.HTTPFaultInjection_Abort
	if fault.Abort != nil {
		abort = &v1alpha3.HTTPFaultInjection_Abort{
			Percentage: convertPercentage(fault.Abort.Percentage),
		}
		if fault.Abort.ErrorType != nil {
			switch errType := fault.Abort.ErrorType.(type) {
			case *v1.HTTPFaultInjection_Abort_HttpStatus:
				abort.ErrorType = &v1alpha3.HTTPFaultInjection_Abort_HttpStatus{
					HttpStatus: errType.HttpStatus,
				}
			case *v1.HTTPFaultInjection_Abort_GrpcStatus:
				abort.ErrorType = &v1alpha3.HTTPFaultInjection_Abort_GrpcStatus{
					GrpcStatus: errType.GrpcStatus,
				}
			case *v1.HTTPFaultInjection_Abort_Http2Error:
				abort.ErrorType = &v1alpha3.HTTPFaultInjection_Abort_Http2Error{
					Http2Error: errType.Http2Error,
				}
			}
		}
	}
	return &v1alpha3.HTTPFaultInjection{
		Delay: delay,
		Abort: abort,
	}
}

func convertPercentage(percent *v1.Percent) *v1alpha3.Percent {
	if percent == nil {
		return nil
	}
	return &v1alpha3.Percent{
		Value: percent.Value,
	}
}

func convertCorsPolicy(cors *v1.CorsPolicy) *v1alpha3.CorsPolicy {
	if cors == nil {
		return nil
	}
	return &v1alpha3.CorsPolicy{
		AllowOrigin:      cors.AllowOrigin,
		AllowMethods:     cors.AllowMethods,
		AllowHeaders:     cors.AllowHeaders,
		ExposeHeaders:    cors.ExposeHeaders,
		MaxAge:           cors.MaxAge,
		AllowCredentials: cors.AllowCredentials,
	}
}

func convertMatch(match []*gloov1.Matcher) []*v1alpha3.HTTPMatchRequest {
	var istioMatch []*v1alpha3.HTTPMatchRequest
	for _, m := range match {
		istioMatch = append(istioMatch, &v1alpha3.HTTPMatchRequest{
			Uri:     convertStringMatch(m.Uri),
			Method:  convertStringMatch(m.Method),
			Headers: convertHeaders(m.Headers),
			// TODO: port and sourcelabels?
		})
	}
	return istioMatch
}

func convertAction(route *gloov1.RouteAction, upstreams gloov1.UpstreamList) ([]*v1alpha3.HTTPRouteDestination, error) {
	switch dest := route.Destination.(type) {
	case *gloov1.RouteAction_Single:
		istioDestination, err := convertDestination(route, upstreams)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert destination %v", destination)
		}
		return []*v1alpha3.HTTPRouteDestination{{
			Destination: istioDestination,
			// TODO: ilackarms: add support in our api
			//RemoveRequestHeaders:  route.RemoveRequestHeaders,
			//RemoveResponseHeaders: route.RemoveResponseHeaders,
			//AppendRequestHeaders:  route.AppendRequestHeaders,
			//AppendResponseHeaders: route.AppendResponseHeaders,
		}}, nil

	case *gloov1.RouteAction_Multi:
	}
}

func convertDestination(ctx context.Context, dest *gloov1.Destination, upstreams gloov1.UpstreamList) (*v1alpha3.Destination, error) {
	upstream, err := upstreams.Find(dest.Upstream.Namespace, dest.Upstream.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid destination %v", dest)
	}
	hosts, err := getHostsForUpstream(upstream)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid upstream %v", upstream)
	}
	if len(hosts) < 1 {
		return nil, errors.Errorf("could not find at least 1 host for upstream %v", upstream)
	}

	var portSelector *v1alpha3.PortSelector
	port, err := getPortForUpstream(upstream)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("no port found for %v, assuming it listeons on only 1 port")
	} else {
		portSelector = &v1alpha3.PortSelector{
			Port: &v1alpha3.PortSelector_Number{
				Number: port,
			},
		}
	}

	return &v1alpha3.Destination{
		Host:   hosts[0], // ilackarms: this host must match what istio expects in the service registry
		Subset: subsetName(upstream),
		Port:   portSelector,
	}, nil
}

func convertHeaders(headers map[string]*v1.StringMatch) map[string]*v1alpha3.StringMatch {
	out := make(map[string]*v1alpha3.StringMatch)
	for k, v := range headers {
		out[k] = convertStringMatch(v)
	}
	return out
}

func convertStringMatch(match *v1.StringMatch) *v1alpha3.StringMatch {
	switch strMatch := match.MatchType.(type) {
	case *v1.StringMatch_Exact:
		return &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Exact{Exact: strMatch.Exact}}
	case *v1.StringMatch_Prefix:
		return &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Prefix{Prefix: strMatch.Prefix}}
	case *v1.StringMatch_Regex:
		return &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Regex{Regex: strMatch.Regex}}
	}
	return nil
}
