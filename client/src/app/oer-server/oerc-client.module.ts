/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2024 emschu[aet]mailbox.org
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
import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {StatusComponent} from './status/status.component';
import {ApiService} from './api.service';
import {TimelineComponent} from './timeline/timeline.component';
import {BrowserModule} from '@angular/platform-browser';
import {HttpClientModule} from '@angular/common/http';
import {RecommendationComponent} from './recommendation/recommendation.component';
import {DashboardComponent} from './dashboard/dashboard.component';
import {CreditsComponent} from './credits/credits.component';
import {SearchComponent} from './search/search.component';
import {LogDashboardComponent} from './log-dashboard/log-dashboard.component';
import {ReadMorePipe} from './read-more.pipe';
import {SearchPipe} from './search.pipe';
import {SearchService} from './search/search.service';
import {UtilModule} from '../util/util.module';
import {FormsModule} from '@angular/forms';
import {StateService} from './state.service';
import {XmltvComponent} from "./xmltv/xmltv.component";

@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
    BrowserModule,
    UtilModule,
    FormsModule,
  ],
  declarations: [
    StatusComponent,
    TimelineComponent,
    RecommendationComponent,
    DashboardComponent,
    XmltvComponent,
    CreditsComponent,
    SearchComponent,
    LogDashboardComponent,
    ReadMorePipe,
    SearchPipe,
  ],
  exports: [
    StatusComponent,
    TimelineComponent,
    RecommendationComponent,
    DashboardComponent,
    LogDashboardComponent,
    CreditsComponent,
    SearchComponent,
  ],
  providers: [
    ApiService,
    SearchService,
    StateService,
  ]
})
export class OercClientModule {
}

