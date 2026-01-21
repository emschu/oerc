/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
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

import {Component, inject, OnInit, signal} from '@angular/core';
import {first} from "rxjs/operators";
import {Channel, ChannelResponse} from "../entities";
import {ApiService} from "../api.service";
import {CdkDrag, CdkDragDrop, CdkDropList, moveItemInArray} from "@angular/cdk/drag-drop";

@Component({
  selector: 'app-channels',
  imports: [CdkDrag, CdkDropList],
  templateUrl: './channels.component.html',
  styleUrl: './channels.component.scss',
})
export class ChannelsComponent implements OnInit {

  channels = signal<Channel[]>([]);

  private apiService = inject(ApiService);

  ngOnInit(): void {
    this.apiService.channels().pipe(first()).subscribe((value: ChannelResponse) => {
      this.channels.set(value.data);
    });
  }

  protected drop($event: CdkDragDrop<Channel[]>) {
    const updatedChannels = [...this.channels()];
    moveItemInArray(updatedChannels, $event.previousIndex, $event.currentIndex);
    this.channels.set(updatedChannels);
    this.apiService.updateChannelsOrder(updatedChannels).pipe(first()).subscribe();
  }
}
