/*
 * oerc, alias oer-collector
 * Copyright (C) 2021 emschu[aet]mailbox.org
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
import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {IdType, Timeline, TimelineEventPropertiesResult, TimelineOptions, TimelineWindow} from 'vis-timeline';
import {Channel, ChannelResponse, ProgramEntry} from '../entities';
import {DataSet} from 'vis-data';
import {DeepPartial} from 'vis-data/declarations/data-interface';
import {Subscription} from 'rxjs';
import moment, {MomentInput} from 'moment-timezone';
import {environment} from '../../../environments/environment';

moment.locale('de');

@Component({
  selector: 'app-oer-timeline',
  templateUrl: './timeline.component.html',
  styleUrls: ['./timeline.component.scss']
})
export class TimelineComponent implements OnInit, OnDestroy {
  public timeLine: Timeline | null = null;
  public items: DataSet<any>;
  currentProgramEntry: ProgramEntry | null = null;
  isModalOpen = false;

  channelSubscription: Subscription | null = null;
  programSubscription: Subscription | null = null;
  loadingSubscription: Subscription | null = null;

  constructor(public apiService: ApiService) {
    this.items = new DataSet<any>();
  }

  ngOnInit(): void {
    this.initTimeLine();
    setTimeout(() => this.moveToNow(), 0);
  }

  ngOnDestroy(): void {
    this.channelSubscription?.unsubscribe();
    this.loadingSubscription?.unsubscribe();
    this.programSubscription?.unsubscribe();
    this.timeLine?.destroy();
    this.items.clear();
  }

  initTimeLine(): void {
    // DOM element where the Timeline will be attached
    const container = document.getElementById('program_timeline');
    if (!container) {
      console.error('Missing element #timeline');
      return;
    }
    const today = moment().tz(environment.timezone).toDate();

    this.programSubscription = this.apiService.programSubject.subscribe(programResponse => {
      if (!programResponse) {
        return;
      }

      const timeZoneOffset = moment(new Date()).tz(environment.timezone).utcOffset();
      const programList: DeepPartial<any> = [];
      programResponse.program_list.forEach(singleProgramEntry => {
        programList.push({
          id: singleProgramEntry.id,
          group: singleProgramEntry.channel_id,
          start: moment.parseZone(singleProgramEntry.start_date_time).subtract(timeZoneOffset, 'minutes'),
          end: moment.parseZone(singleProgramEntry.end_date_time).subtract(timeZoneOffset, 'minutes'),
          content: singleProgramEntry.title,
          title: singleProgramEntry.title,
          type: 'range',
        });
      });
      this.items.update(programList);
    });

    this.apiService.fetchProgramForDay(today);

    // create groups
    const groups: DataSet<any> = new DataSet({});
    this.channelSubscription = this.apiService.channels().subscribe((value: ChannelResponse) => {
      if (!value) {
        return;
      }
      value.data.forEach(
        (singleChannel: Channel) => {
          groups.add({id: singleChannel.id, content: singleChannel.title});
        });
    });

    // Configuration for the Timeline
    const options: TimelineOptions = {
      stack: false,
      start: moment().tz(environment.timezone).toDate(),
      end: moment().tz(environment.timezone).add(2, 'hours').toDate(),
      editable: false,
      orientation: 'top',
      zoomable: true,
      showCurrentTime: true,
      clickToUse: false,
      horizontalScroll: false,
      verticalScroll: true,
      zoomKey: 'ctrlKey',
      zoomMax: 300000000,
      maxHeight: 550,
      minHeight: 450,
      moveable: false,
      multiselect: false,
      multiselectPerGroup: false,
      rtl: false,
      selectable: true,
      moment: (date: MomentInput) => {
        return moment.utc(date).tz(environment.timezone, false);
      },
      margin: {
        item: 5,
        axis: 1
      },
      tooltip: {
        followMouse: true,
        overflowMethod: `cap`
      },
    };

    // Create a Timeline
    this.timeLine = new Timeline(container, this.items, groups, options);

    this.timeLine.on('rangechanged', this.rangeChange.bind(this));
    this.timeLine.on('doubleClick', this.itemClicked.bind(this));

    // TODO replace by onClick-handler
    const zoomInElement = document.getElementById('zoomIn');
    const zoomOutElement = document.getElementById('zoomOut');
    const moveLeftElement = document.getElementById('moveLeft');
    const moveRightElement = document.getElementById('moveRight');
    const nowTimelineElement = document.getElementById('nowTimeline');
    if (zoomInElement) {
      zoomInElement.onclick = this.zoomIn.bind(this);
    }
    if (zoomOutElement) {
      zoomOutElement.onclick = this.zoomOut.bind(this);
    }
    if (moveLeftElement) {
      moveLeftElement.onclick = this.moveLeft.bind(this);
    }
    if (moveRightElement) {
      moveRightElement.onclick = this.moveRight.bind(this);
    }
    if (nowTimelineElement) {
      nowTimelineElement.onclick = this.moveToNow.bind(this);
    }
  }

  zoomIn(): void {
    this.timeLine?.zoomIn(0.25);
  }

  zoomOut(): void {
    this.timeLine?.zoomOut(0.25);
  }

  moveLeft(): void {
    this.move(0.25);
  }

  moveRight(e: Event): void {
    this.move(-0.25);
  }

  move(percentage: number): void {
    if (!this.timeLine) {
      return;
    }

    const range: TimelineWindow = this.timeLine.getWindow();
    const interval: number = range.end.valueOf() - range.start.valueOf();

    this.timeLine.setWindow(
      range.start.valueOf() - interval * percentage,
      range.end.valueOf() - interval * percentage
    );
  }

  private rangeChange(e: Event): void {
    if (this.timeLine === undefined) {
      return;
    }
    const rangeStart = this.timeLine?.getWindow().start;
    const rangeEnd = this.timeLine?.getWindow().end;

    if (!rangeStart || !rangeEnd) {
      return;
    }

    if (!this.apiService.checkIfDayIsFetched(rangeStart)) {
      this.apiService.fetchProgramForDay(new Date(rangeStart.getFullYear(), rangeStart.getMonth(), rangeStart.getDate()));
    } else if (!this.apiService.checkIfDayIsFetched(rangeEnd)) {
      this.apiService.fetchProgramForDay(new Date(rangeEnd.getFullYear(), rangeEnd.getMonth(), rangeEnd.getDate()));
    }
  }

  itemClicked(e: TimelineEventPropertiesResult): void {
    if (e.item) {
      const clickedEntryId: IdType | null = e.item;
      if (!clickedEntryId) {
        return;
      }
      this.loadingSubscription = this.apiService.entry(clickedEntryId).subscribe(value => {
        this.isModalOpen = true;
        this.currentProgramEntry = value;
      });
    }
  }

  @HostListener('document:keyup', ['$event'])
  handleKeyboardEvent(event: KeyboardEvent): void {
    // @ts-ignore
    if (event.target?.nodeName?.toUpperCase() === 'INPUT') {
      // ignore key events on input elements
      return;
    }

    if (event.key === 'Escape') {
      this.isModalOpen = false;
    }
    if (event.key === 'r' || event.key === 'ArrowRight') {
      this.moveRight(event);
    }
    if (event.key === 'l' || event.key === 'ArrowLeft') {
      this.moveLeft();
    }
    if (event.key === 'i') {
      this.zoomIn();
    }
    if (event.key === 'o') {
      this.zoomOut();
    }
  }

  moveToNow(): void {
    this.timeLine?.moveTo(new Date(), {animation: true});
  }
}
