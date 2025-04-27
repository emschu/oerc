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
import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {RecommendationComponent} from './oer-server/recommendation/recommendation.component';
import {DashboardComponent} from './oer-server/dashboard/dashboard.component';
import {NotFoundComponent} from './not-found/not-found.component';
import {CreditsComponent} from './oer-server/credits/credits.component';
import {SearchComponent} from './oer-server/search/search.component';
import {LogDashboardComponent} from './oer-server/log-dashboard/log-dashboard.component';
import {XmltvComponent} from "./oer-server/xmltv/xmltv.component";

const routes: Routes = [
  {path: '', component: DashboardComponent},
  {path: 'recommendations', component: RecommendationComponent},
  {path: 'xmltv', component: XmltvComponent},
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
