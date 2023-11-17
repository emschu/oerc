/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2023 emschu[aet]mailbox.org
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
import {AfterViewInit, Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {DataGroup, IdType, Timeline, TimelineEventPropertiesResult, TimelineOptions, TimelineWindow} from 'vis-timeline/esnext/esm';
import {Channel, ChannelResponse, ProgramEntry, ProgramEntryEssential} from '../entities';
import {BehaviorSubject, Subscription} from 'rxjs';
import {environment} from '../../../environments/environment';
import {first, skip} from 'rxjs/operators';
import {StateService} from '../state.service';
import flatpickr from 'flatpickr';
import * as flatPickrLang from 'flatpickr/dist/l10n/de';
import * as visDataTypes from 'vis-data/declarations/data-interface';
import {UpdateItem} from 'vis-data/declarations/data-interface';
import {DataSet} from 'vis-data/esnext/esm';
import dayjs from 'dayjs';
import {DataItem} from 'vis-timeline';
import {DataInterface} from 'vis-data';
import FlatPickrInstance = flatpickr.Instance;

// interface GroupOrder {
//   groupId: number | string;
//   title: string;
//   order: number;
// }

@Component({
    selector: 'app-oer-timeline',
    templateUrl: './timeline.component.html',
    styleUrls: ['./timeline.component.scss']
})
export class TimelineComponent implements OnInit, OnDestroy, AfterViewInit {

    constructor(public apiService: ApiService,
                private stateService: StateService) {
        this.items = new DataSet();
    }

    private static i = 0;
    static channelMap = new Map<number, number>([
        [1, TimelineComponent.i++],
        [16, TimelineComponent.i++],
        [9, TimelineComponent.i++],
        [5, TimelineComponent.i++],
        [20, TimelineComponent.i++],
        [21, TimelineComponent.i++],
        [17, TimelineComponent.i++],
        [18, TimelineComponent.i++],
        [14, TimelineComponent.i++],
        [13, TimelineComponent.i++],
        [2, TimelineComponent.i++],
        [12, TimelineComponent.i++],
        [11, TimelineComponent.i++],
        [4, TimelineComponent.i++],
        [19, TimelineComponent.i++],
    ]);
    public items: DataInterface<DataItem, 'id'>;
    public timeLine?: Timeline;
    currentProgramEntry?: ProgramEntry;
    isModalOpen = false;
    isMissingDataModalOpen = false;

    // bound to form-switch-control + initial value
    showDeprecatedEntries = new BehaviorSubject(this.stateService.getShowDeprecatedEntries());

    // subscriptions managed by this component
    channelSubscription?: Subscription;
    programSubscription?: Subscription;
    loadingSubscription?: Subscription;
    showDeprecatedEntriesSubscription?: Subscription;
    // private latestSelectedDate: Date | null = null;
    private dateTimePickrInstance?: FlatPickrInstance;

    private readonly _datePickerFormat = 'DD.MM.YY HH:mm';

    /**
     * TODO: make configuration option
     * @param singleChannelID an id
     * @private
     */
    private static getGroupOrder(singleChannelID: number): number {
        if (TimelineComponent.channelMap.has(singleChannelID)) {
            return TimelineComponent.channelMap.get(singleChannelID) as number;
        }
        return 100;
    }

    ngOnInit(): void {
        this.initTimeLine();

        this.apiService.statusSubject.pipe(first()).subscribe(statusResponse => {
            this.dateTimePickrInstance = flatpickr('#timeline_date_range_picker', {
                locale: flatPickrLang.German,
                now: dayjs().locale(environment.locale).format(),
                enableTime: true,
                allowInput: false,
                time_24hr: true,
                clickOpens: true,
                dateFormat: this._datePickerFormat,
                altFormat: this._datePickerFormat,
                defaultHour: 18,
                enableSeconds: false,
                minuteIncrement: 15,
                mode: 'single',
                defaultDate: dayjs().locale(environment.locale).format(),
                onChange: (selectedDates: Date[], dateStr: string, _: FlatPickrInstance) => {
                    if (selectedDates.length === 0) {
                        return;
                    }
                    if (this.timeLine) {
                        this.timeLine.moveTo(dayjs(selectedDates[0]).toISOString(), {animation: false});
                    }
                },
                parseDate: (dateString, format) => {
                    const timezonedDate = dayjs(dateString, format).locale(environment.locale);
                    return new Date(
                        timezonedDate.year(),
                        timezonedDate.month(),
                        timezonedDate.date(),
                        timezonedDate.hour(),
                        timezonedDate.minute()
                    );
                },
                formatDate: (date, format) => {
                    return dayjs(date).format(format);
                },
            }) as FlatPickrInstance;

            if (statusResponse?.data_start_time && statusResponse?.data_end_time) {
                this.dateTimePickrInstance.set({
                    minDate: dayjs(statusResponse?.data_start_time).format(),
                    maxDate: dayjs(statusResponse?.data_end_time).format(),
                });
            }
        });
        this.moveToNow();
    }

    ngOnDestroy(): void {
        this.channelSubscription?.unsubscribe();
        this.loadingSubscription?.unsubscribe();
        this.programSubscription?.unsubscribe();
        this.showDeprecatedEntriesSubscription?.unsubscribe();
        this.showDeprecatedEntries.unsubscribe();

        this.timeLine?.destroy();
    }

    ngAfterViewInit(): void {
        this.moveToNow();
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
        const groups: DataSet<DataGroup> = new DataSet({fieldId: 'id'});
        this.channelSubscription = this.apiService.channels().pipe(first()).subscribe((value: ChannelResponse) => {
            if (!value) {
                return;
            }
            value.data.forEach(
                (singleChannel: Channel) => {
                    groups.add({
                        id: singleChannel.id,
                        content: singleChannel.title,
                        subgroupStack: true,
                        subgroupOrder: () => 0,
                    });
                });
        });

        // Configuration for the Timeline
        const now = dayjs().locale(environment.locale);
        const options: TimelineOptions = {
            align: 'center',
            locale: environment.locale,
            stack: false,
            stackSubgroups: true,
            start: now.clone().subtract(1, 'hour').toISOString(),
            end: now.clone().add(3, 'hour').toISOString(),
            timeAxis: {scale: 'minute', step: 15},
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
            groupOrder(a: DataGroup, b: DataGroup): number {
                if (a.id === b.id) {
                    return 0;
                }
                return TimelineComponent.getGroupOrder(a.id as number) > TimelineComponent.getGroupOrder(b.id as number) ? 1 : -1;
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

        this.showDeprecatedEntriesSubscription = this.showDeprecatedEntries.pipe(skip(1)).subscribe(value => {
            this.stateService.setShowDeprecatedEntries(value);

            if (!value) {
                this.items = new DataSet();
            }
            this.loadProgramItems();
        });
    }

    private loadProgramItems(): void {
        const now = dayjs().locale(environment.locale);
        const midnight = dayjs().locale(environment.locale).hour(0).minute(0).second(0);

        if (!this.programSubscription || this.programSubscription.closed) {
            this.programSubscription = this.apiService.programSubject.subscribe(programResponse => {
                if (!programResponse || programResponse.program_list?.length === 0) {
                    return;
                }
                this.apiService.isLoadingSubject.next(true);

                const showDeprecatedEntries = this.showDeprecatedEntries.getValue();
                const programEntries = programResponse.program_list.flatMap((value): ProgramEntryEssential => {
                    return {
                        id: value.id,
                        created_at: value.created_at,
                        start_date_time: value.start_date_time,
                        end_date_time: value.end_date_time,
                        channel_id: value.channel_id,
                        is_deprecated: value.is_deprecated,
                        title: value.title,
                        hash: value.hash,
                    };
                });

                const programList: visDataTypes.DeepPartial<DataItem[]> = [];

                function getAdditionalTitleInfo(singleProgramEntry: ProgramEntryEssential): string {
                    // todo i18n
                    return ' | CreatedAt: ' + dayjs(singleProgramEntry.created_at).locale(environment.locale).format('D.M HH:mm:ss');
                }

                programEntries.forEach(singleProgramEntry => {
                    if (!showDeprecatedEntries && singleProgramEntry.is_deprecated) {
                        // just ignore them
                        return;
                    }
                    programList.push({
                        id: singleProgramEntry.id,
                        group: singleProgramEntry.channel_id,
                        start: singleProgramEntry.start_date_time,
                        end: singleProgramEntry.end_date_time,
                        content: singleProgramEntry.title,
                        title: singleProgramEntry.title + getAdditionalTitleInfo(singleProgramEntry),
                        type: 'range',
                        subgroup: 1,
                        className: singleProgramEntry.is_deprecated ? 'deprecated-item' : '',
                    });
                });
                this.items.getDataSet().update(programList);

                if (showDeprecatedEntries) {
                    const deprecatedEntries: visDataTypes.DeepPartial<UpdateItem<DataItem, 'id'>[]> = [];
                    programEntries.filter(value => value.is_deprecated).forEach(singleProgramEntry => {
                        // this is a very expensive loop
                        const overlaps = this.items.get({
                            filter: item => {
                                if (item.group !== singleProgramEntry.channel_id) {
                                    return false;
                                }
                                if (item.id === singleProgramEntry.id) {
                                    return false;
                                }
                                if (item.className !== 'deprecated-item') {
                                    return false;
                                }
                                const t1 = singleProgramEntry.start_date_time;
                                const t2 = singleProgramEntry.end_date_time;
                                const t3 = item.start as Date;
                                const t4 = item.end as Date;

                                if ((t1 <= t3 && t2 <= t3)
                                    || (t1 > t4)
                                    || (t2 === t3)
                                    || (t1 === t4)
                                ) {
                                    return false;
                                }
                                // entries could overlap in 9 ways:
                                if ((t1 === t3 && t2 === t4)
                                    || (t1 === t3 && t2 < t4)
                                    || (t1 === t3 && t2 > t4)
                                    || (t1 > t3 && t2 === t4)
                                    || (t1 < t3 && t2 === t4)
                                    || (t1 > t3 && t2 < t4)
                                    || (t1 < t3 && t2 > t4)
                                    || (t1 > t3 && t2 > t4)
                                    || (t1 < t3 && t2 < t4)
                                ) {
                                    return true;
                                }
                                return false;
                            }
                        }).sort((a, b) => a.id < b.id ? 1 : -1);

                        // this mechanism tries to find direct neighbours to exclude them from subgroup calculation
                        const excludedOverlaps: Set<IdType> = new Set<IdType>();
                        overlaps.forEach(singleItem => {
                            overlaps.forEach(otherItem => {
                                if (singleItem === otherItem) {
                                    return;
                                }
                                if (singleItem.end === otherItem.start) {
                                    excludedOverlaps.add(singleItem.id);
                                    return;
                                }
                            });
                        });

                        let subgroupID = 2;
                        if (overlaps.length > 0) {
                            const affectedIDs = overlaps
                                .filter(item => !excludedOverlaps.has(item.id))
                                .flatMap(value => value.id);
                            affectedIDs.push(singleProgramEntry.id);
                            const newIndex = affectedIDs.sort((a, b) => {
                                if (a === b) {
                                    return 0;
                                }
                                return a < b ? 1 : -1;
                            }).findIndex(value => value === singleProgramEntry.id);
                            subgroupID += newIndex;
                        }

                        deprecatedEntries.push({
                            id: singleProgramEntry.id,
                            group: singleProgramEntry.channel_id,
                            start: singleProgramEntry.start_date_time,
                            end: singleProgramEntry.end_date_time,
                            content: singleProgramEntry.title,
                            title: singleProgramEntry.title + getAdditionalTitleInfo(singleProgramEntry),
                            type: 'range',
                            subgroup: subgroupID,
                            className: 'deprecated-item'
                        });
                    });
                    this.items.getDataSet().updateOnly(deprecatedEntries);
                }

                setTimeout(() => {
                    this.apiService.isLoadingSubject.next(false);
                }, 500);
            });
        }

        let range = this.timeLine?.getWindow();
        let currentTlTime;
        if (range) {
            currentTlTime = dayjs(range.start.valueOf()).locale(environment.locale);
        } else {
            currentTlTime = now;
        }

        // load yesterday's program if we are just after midnight to show enough items in the timeline
        const minuteDiff = now.diff(midnight, 'minute', false);
        this.apiService.fetchProgramForDay(currentTlTime.toDate());
        if (minuteDiff < 180) {
            this.apiService.fetchProgramForDay(now.clone().subtract(1, 'day').toDate());
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
        this.dateTimePickrInstance?.setDate(dayjs().locale(environment.locale).format(), false);
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

        this.apiService.fetchProgramForDay(new Date(rangeStart.getFullYear(), rangeStart.getMonth(), rangeStart.getDate()));
        if (rangeStart.getDay() != rangeEnd.getDay()) {
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
