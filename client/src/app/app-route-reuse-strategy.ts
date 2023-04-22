import {ActivatedRouteSnapshot, DetachedRouteHandle, RouteReuseStrategy} from '@angular/router';

export class AppRouteReuseStrategy implements RouteReuseStrategy {
  cachedRoutes: string[] = ['', 'recommendations', 'log', 'credits'];
  samePageRefreshRoutes: string[] = ['search'];
  routeHandles = new Map<string, DetachedRouteHandle>();

  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    return this.cachedRoutes.indexOf(route.routeConfig?.path ?? '') > -1;
  }

  store(route: ActivatedRouteSnapshot, handle: DetachedRouteHandle): void {
    this.routeHandles.set(route.routeConfig?.path ?? '', handle);
  }

  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    return this.routeHandles.has(route.routeConfig?.path ?? '');
  }

  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle {
    // @ts-ignore
    return this.routeHandles.get(route.routeConfig?.path ?? '');
  }

  shouldReuseRoute(future: ActivatedRouteSnapshot, current: ActivatedRouteSnapshot): boolean {
    if (this.samePageRefreshRoutes.indexOf(current.routeConfig?.path ?? '') > -1) {
      return false;
    }
    return future.routeConfig === current.routeConfig;
  }
}
