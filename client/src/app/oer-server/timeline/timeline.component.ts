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
import {Channel, ChannelResponse, ProgramEntry, ProgramResponse} from '../entities';
import {DataSet} from 'vis-data';
import {DeepPartial} from 'vis-data/declarations/data-interface';
import {BehaviorSubject, Subscription} from 'rxjs';
import moment, {MomentInput} from 'moment-timezone';
import {environment} from '../../../environments/environment';
import {first, skip, take} from 'rxjs/operators';
import {StateService} from '../state.service';
import flatpickr from 'flatpickr';
import FlatPickrInstance = flatpickr.Instance;

interface SubgroupMappingChannelLevel {
  [key: number]: ProgramEntry[];
}

interface SubgroupMappingEntryLevel {
  [key: string]: number;
}

interface GroupOrder {
  groupId: number | string;
  title: string;
  order: number;
}

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

  // bound to form-switch-control + initial value
  showDeprecatedEntries = new BehaviorSubject(this.stateService.getShowDeprecatedEntries());

  // subscriptions managed by this component
  channelSubscription: Subscription | null = null;
  programSubscription: Subscription | null = null;
  loadingSubscription: Subscription | null = null;
  showDeprecatedEntriesSubscription: Subscription | null = null;
  // private latestSelectedDate: Date | null = null;
  private dateTimePickrInstance: FlatPickrInstance | null = null;
  private zoneOffset = 0;

  constructor(public apiService: ApiService,
              private stateService: StateService) {
    this.items = new DataSet<any>();
  }

  ngOnInit(): void {
    this.zoneOffset = moment().tz(environment.timezone).utcOffset();
    this.initTimeLine();
    setTimeout(() => this.moveToNow(), 0);

    this.apiService.statusSubject.pipe(skip(1), take(1)).subscribe(value => {
      this.dateTimePickrInstance = flatpickr('#timeline_date_range_picker', {
        now: moment().tz(environment.timezone).toISOString(),
        enableTime: true,
        allowInput: true,
        time_24hr: true,
        clickOpens: true,
        altFormat: 'd.m.Y H:m',
        defaultHour: 18,
        enableSeconds: false,
        minuteIncrement: 15,
        mode: 'single',
        defaultDate: moment().tz(environment.timezone).toISOString(),
        minDate: moment(value?.data_start_time).tz(environment.timezone).utcOffset(this.zoneOffset).format(),
        maxDate: moment(value?.data_end_time).tz(environment.timezone).utcOffset(this.zoneOffset).format(),
        onChange: (selectedDates: Date[], dateStr: string, instance: FlatPickrInstance) => {
          if (selectedDates.length === 0) {
            return;
          }
          this.timeLine?.moveTo(selectedDates[0].toISOString(), {animation: true});
        },
      }) as FlatPickrInstance;
    });
    this.apiService.updateStatus();
  }

  ngOnDestroy(): void {
    this.channelSubscription?.unsubscribe();
    this.loadingSubscription?.unsubscribe();
    this.programSubscription?.unsubscribe();
    this.showDeprecatedEntriesSubscription?.unsubscribe();

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
    this.loadProgramItems();

    // create groups
    const groups: DataSet<any> = new DataSet({});
    this.channelSubscription = this.apiService.channels().pipe(first()).subscribe((value: ChannelResponse) => {
      if (!value) {
        return;
      }
      value.data.forEach(
        (singleChannel: Channel) => {
          groups.add({
            id: singleChannel.id,
            content: singleChannel.title,
            editable: false,
            subgroupStack: false,
            order: this.getGroupOrder(singleChannel),
            subgroupOrder: () => 0,
          });
        });
    });

    // Configuration for the Timeline
    const now = moment().tz(environment.timezone).utcOffset(this.zoneOffset);
    const options: TimelineOptions = {
      align: 'center',
      locale: 'de',
      stack: false,
      stackSubgroups: false,
      start: now.clone().subtract(3, 'hours').format(),
      end: now.clone().add(3, 'hours').format(),
      timeAxis: {scale: 'minute', step: 15},
      moment: (date: MomentInput | undefined) => {
        return moment(date).tz(environment.timezone).utcOffset(this.zoneOffset);
      },
      orientation: 'top',
      zoomable: true,
      showCurrentTime: true,
      clickToUse: false,
      horizontalScroll: false,
      verticalScroll: true,
      zoomMin: 1000 * 60 * 60 * 2,
      zoomKey: 'ctrlKey',
      zoomMax: 200000000,
      maxHeight: 550,
      minHeight: 450,
      moveable: false,
      multiselect: false,
      multiselectPerGroup: false,
      rtl: false,
      selectable: true,
      editable: false,
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

    this.showDeprecatedEntriesSubscription = this.showDeprecatedEntries.pipe(skip(1)).subscribe(value => {
      this.stateService.setShowDeprecatedEntries(value);

      this.items.clear();
      this.loadProgramItems();
    });
  }

  /**
   * TODO: make configuration option
   * @param singleChannel
   * @private
   */
  private getGroupOrder(singleChannel: Channel): number {
    let i = 0;
    const channelMap = new Map<number, number>([
      [1, i++],
      [16, i++],
      [9, i++],
      [5, i++],
      [20, i++],
      [21, i++],
      [17, i++],
      [18, i++],
      [14, i++],
      [13, i++],
      [2, i++],
      [12, i++],
      [11, i++],
      [4, i++],
      [19, i++],
    ]);
    if (channelMap.has(singleChannel.id)) {
      return channelMap.get(singleChannel.id) as number;
    }
    return 100;
  }

  private loadProgramItems(): void {
    this.programSubscription?.unsubscribe();

    const today = moment().tz(environment.timezone).toDate();

    this.programSubscription = this.apiService.programSubject.subscribe(programResponse => {
      if (!programResponse) {
        return;
      }
      this.apiService.isLoadingSubject.next(true);

      let subgroupMappingOfEntries: SubgroupMappingEntryLevel = {};
      if (this.showDeprecatedEntries.getValue()) {
        // this is pretty costly and memory intensive...
        subgroupMappingOfEntries = this.calculateSubgroupMapping(programResponse);
      }

      const programList: DeepPartial<any> = [];
      programResponse.program_list.forEach(singleProgramEntry => {
        if (singleProgramEntry.is_deprecated && !this.showDeprecatedEntries.getValue()) {
          // skip deprecated entries, if the user wants this
          return;
        }
        let subGroupId = 1; // default
        let defaultClass = '';
        if (singleProgramEntry.is_deprecated) {
          if (subgroupMappingOfEntries[singleProgramEntry.id] === undefined) {
            console.error('unknown subgroup id for entry with id #' + singleProgramEntry.id);
          }
          subGroupId = subgroupMappingOfEntries[singleProgramEntry.id] + 1;
          defaultClass = 'deprecated-item';
        }
        programList.push({
          id: singleProgramEntry.id,
          group: singleProgramEntry.channel_id,
          start: moment(singleProgramEntry.start_date_time).tz(environment.timezone).utcOffset(this.zoneOffset),
          end: moment(singleProgramEntry.end_date_time).tz(environment.timezone).utcOffset(this.zoneOffset),
          content: singleProgramEntry.title,
          title: singleProgramEntry.title,
          type: 'range',
          subgroup: subGroupId,
          className: defaultClass
        });
      });
      this.items.update(programList);
      setTimeout(() => {
        this.apiService.isLoadingSubject.next(false);
      }, 500);
    });

    this.apiService.fetchProgramForDay(today);
  }

  private calculateSubgroupMapping(programResponse: ProgramResponse): SubgroupMappingEntryLevel {
    // helper function
    const groupBy = (xs: any, key: any): {} => {
      return xs.reduce((rv: any, x: any) => {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
      }, {});
    };
    const subgroupChannelMap: SubgroupMappingChannelLevel = [];
    const deprecatedEntries = groupBy(programResponse.program_list.filter(item => item.is_deprecated), 'channel_id');
    Object.entries(deprecatedEntries).forEach((value) => {
      subgroupChannelMap[parseInt(value[0], 10)] = (value[1] as ProgramEntry[]).sort((a, b) => {
        return moment(b.start_date_time).unix() - moment(a.start_date_time).unix();
      });
    });
    const subgroupEntryMap: SubgroupMappingEntryLevel = {};

    for (const key of Object.keys(deprecatedEntries)) {
      const channelId = parseInt(key, 10);
      const deprecatedEntriesCountInChannel = subgroupChannelMap[channelId].length;
      for (let i = 0; i < deprecatedEntriesCountInChannel; i++) {
        if (!subgroupEntryMap[subgroupChannelMap[channelId][i].id]) {
          subgroupEntryMap[subgroupChannelMap[channelId][i].id] = 1;
        } else {
          subgroupEntryMap[subgroupChannelMap[channelId][i].id] += 1;
        }
        if (i + 1 < deprecatedEntriesCountInChannel) {
          if (subgroupChannelMap[channelId][i].end_date_time > subgroupChannelMap[channelId][i + 1].start_date_time) {
            if (subgroupEntryMap[subgroupChannelMap[channelId][i].id] > 1) {
              subgroupEntryMap[subgroupChannelMap[channelId][i + 1].id] = 0;
            } else {
              subgroupEntryMap[subgroupChannelMap[channelId][i + 1].id] = 1;
            }
          }
        }
      }
    }
    return subgroupEntryMap;
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

  moveRight(): void {
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

  moveToNow(): void {
    this.dateTimePickrInstance?.setDate(moment().tz(environment.timezone).utcOffset(this.zoneOffset).format());
    this.dateTimePickrInstance?._debouncedChange();
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
      this.moveRight();
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
}
