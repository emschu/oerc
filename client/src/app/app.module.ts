import {NgModule} from '@angular/core';
import {BrowserModule} from '@angular/platform-browser';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {OercClientModule} from './oer-server/oerc-client.module';
import {NotFoundComponent} from './not-found/not-found.component';
import {UtilModule} from './util/util.module';
import { NavComponent } from './nav/nav.component';
import { AngularPageVisibilityModule } from 'angular-page-visibility';

@NgModule({
  declarations: [
    AppComponent,
    NotFoundComponent,
    NavComponent,
  ],
  imports: [
    BrowserModule,
    AngularPageVisibilityModule,
    AppRoutingModule,
    OercClientModule,
    UtilModule,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
