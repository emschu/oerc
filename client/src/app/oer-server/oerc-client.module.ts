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

@NgModule({
  imports: [
    CommonModule,
    HttpClientModule,
    BrowserModule,
    UtilModule,
  ],
  declarations: [
    StatusComponent,
    TimelineComponent,
    RecommendationComponent,
    DashboardComponent,
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
  ]
})
export class OercClientModule {
}

