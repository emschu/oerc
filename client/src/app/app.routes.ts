import {Routes} from '@angular/router';
import {RecommendationComponent} from './oer-server/recommendation/recommendation.component';
import {DashboardComponent} from './oer-server/dashboard/dashboard.component';
import {NotFoundComponent} from './not-found/not-found.component';
import {CreditsComponent} from './oer-server/credits/credits.component';
import {SearchComponent} from './oer-server/search/search.component';
import {LogDashboardComponent} from './oer-server/log-dashboard/log-dashboard.component';
import {XmltvComponent} from './oer-server/xmltv/xmltv.component';

export const routes: Routes = [
  {path: '', component: DashboardComponent},
  {path: 'channels', loadComponent: () => import('./oer-server/channels/channels.component').then(m => m.ChannelsComponent)},
  {path: 'recommendations', component: RecommendationComponent},
  {path: 'xmltv', component: XmltvComponent},
  {path: 'log', component: LogDashboardComponent},
  {path: 'credits', component: CreditsComponent},
  {path: 'search', component: SearchComponent},
  {path: '404', component: NotFoundComponent},
  {path: '**', redirectTo: '/404'}
];
