import {ApplicationConfig, LOCALE_ID, provideZoneChangeDetection} from '@angular/core';
import {provideRouter, RouteReuseStrategy, withHashLocation} from '@angular/router';
import {provideHttpClient, withInterceptorsFromDi, withXhr} from '@angular/common/http';
import {AppRouteReuseStrategy} from './app-route-reuse-strategy';
import {environment} from '../environments/environment';
import {routes} from './app.routes';
import '@angular/common/locales/global/de';
import '@angular/common/locales/global/en';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({eventCoalescing: true}),
    provideRouter(routes, withHashLocation()),
    provideHttpClient(withXhr(), withInterceptorsFromDi()),
    {provide: LOCALE_ID, useValue: environment.locale},
    {provide: RouteReuseStrategy, useClass: AppRouteReuseStrategy}
  ]
};
