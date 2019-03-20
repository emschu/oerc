import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import {StatusComponent} from "./status/status.component";
import {ApiService} from "./api.service";
import { TimelineComponent } from './timeline/timeline.component';

@NgModule({
  imports: [
    CommonModule,
  ],
  declarations: [
    StatusComponent,
    TimelineComponent
  ],
  exports: [
    StatusComponent,
    TimelineComponent
  ],
  providers: [
    ApiService
  ]
})
export class OerServerModule { }
