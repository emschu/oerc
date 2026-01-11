/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
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
import {LOCALE_ID, NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import '@angular/common/locales/global/de';
import '@angular/common/locales/global/en';
import {BrowserModule} from '@angular/platform-browser';
import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {OercClientModule} from './oer-server/oerc-client.module';
import {NotFoundComponent} from './not-found/not-found.component';
import {UtilModule} from './util/util.module';
import {NavComponent} from './nav/nav.component';
import {environment} from '../environments/environment';
import {RouteReuseStrategy} from '@angular/router';
import {AppRouteReuseStrategy} from './app-route-reuse-strategy';

@NgModule({
  declarations: [
    AppComponent,
    NotFoundComponent,
    NavComponent,
  ],
  imports: [
    BrowserModule,
    CommonModule,
    AppRoutingModule,
    OercClientModule,
    UtilModule,
  ],
  bootstrap: [AppComponent],
  providers: [{provide: LOCALE_ID, useValue: environment.locale}, {
    provide: RouteReuseStrategy, useClass: AppRouteReuseStrategy
  }]
})
export class AppModule {
}
