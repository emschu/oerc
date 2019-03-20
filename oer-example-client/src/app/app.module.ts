import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';
import { StatusComponent } from './oer-server/status/status.component';
import {OerServerModule} from "./oer-server/oer-server.module";
import {HttpClientModule} from "@angular/common/http";


@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    OerServerModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
