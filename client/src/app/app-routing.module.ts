import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {RecommendationComponent} from './oer-server/recommendation/recommendation.component';
import {DashboardComponent} from './oer-server/dashboard/dashboard.component';
import {NotFoundComponent} from './not-found/not-found.component';
import {CreditsComponent} from './oer-server/credits/credits.component';
import {SearchComponent} from './oer-server/search/search.component';
import {LogDashboardComponent} from './oer-server/log-dashboard/log-dashboard.component';

const routes: Routes = [
  {path: '', component: DashboardComponent},
  {path: 'recommendations', component: RecommendationComponent},
  {path: 'log', component: LogDashboardComponent},
  {path: 'credits', component: CreditsComponent},
  {path: 'search', component: SearchComponent},
  {path: '404', component: NotFoundComponent},
  {path: '**', redirectTo: '/404'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {useHash: true})],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
