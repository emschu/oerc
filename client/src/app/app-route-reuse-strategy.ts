/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2025 emschu[aet]mailbox.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 */
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
