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
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {StatusResponse} from '../entities';
import {Subscription} from 'rxjs';

@Component({
    selector: 'app-oer-status-display',
    templateUrl: './status.component.html',
    styleUrls: ['./status.component.scss'],
    standalone: false
})
export class StatusComponent implements OnInit, OnDestroy {
  public currentStatus: StatusResponse | null = null;
  statusSubscription: Subscription | null = null;

  constructor(private oerApiService: ApiService) {}

  ngOnInit(): void {
    this.statusSubscription = this.oerApiService.statusSubject.subscribe(value => {
      this.currentStatus = value;
    });
  }

  ngOnDestroy(): void {
    this.statusSubscription?.unsubscribe();
  }
}
