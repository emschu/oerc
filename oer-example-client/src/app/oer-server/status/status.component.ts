import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from "../api.service";
import {StatusResponse} from "../entities";
import {BackgroundService, BackgroundServiceRegistry} from "../updater";
import {DataSet, Timeline} from "vis";

@Component({
  selector: 'oer-status-display',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.css']
})
export class StatusComponent implements OnInit, OnDestroy {
  public currentStatus: StatusResponse;
  private static service_key_status = "status";

  constructor(private oerApiService: ApiService) {
  }

  ngOnInit() {
    this.updateStatus();
    BackgroundServiceRegistry.registerService(new BackgroundService(StatusComponent.service_key_status, this.updateStatus.bind(this), 5000));
  }

  updateStatus() {
    this.oerApiService.statusResponse().subscribe(statusResponse => {
      console.log(statusResponse);
      this.currentStatus = statusResponse;
    });
  }

  ngOnDestroy(): void {
    BackgroundServiceRegistry.unregisterService(StatusComponent.service_key_status);
  }
}
