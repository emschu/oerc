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
import {Injectable} from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {BehaviorSubject, EMPTY, Observable, Subscription} from 'rxjs';
import {
  Channel, ChannelResponse, LogEntryResponse, Pong, ProgramEntry, ProgramResponse, Recommendation, StatusResponse,
  TvShow
} from './entities';
import {IdType} from 'vis-timeline';
import {catchError, first, tap, timeout} from 'rxjs/operators';
import {environment} from '../../environments/environment';
import dayjs from 'dayjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiEndpoint = `${environment.serverEndpoint}`;

  private _isLiveSubject = new BehaviorSubject<boolean | null>(null);
  private _channelSubject = new BehaviorSubject<ChannelResponse | null>(null);
  private _programSubject = new BehaviorSubject<ProgramResponse | null>(null);
  private _tvShowSubject = new BehaviorSubject<TvShow[] | null>(null);
  private _isLoadingSubject = new BehaviorSubject<boolean>(true);
  private _isInErrorsSubject = new BehaviorSubject<boolean>(false);
  private _isWindowOpenedSubject = new BehaviorSubject<boolean>(true);
  private _statusSubject = new BehaviorSubject<StatusResponse | null>(null);

  private channelStore: Channel[] = [];
  private fetchedDays: Date[] = [];
  private isFetchingChannels = false;

  constructor(public http: HttpClient) {
    this.liveCheck();

    this.fetchChannels();

    this.updateStatus();

    this.updateTvShows();
  }

  public init(): void {
    setInterval(() => {
      if (this._isWindowOpenedSubject.getValue()) {
        this.liveCheck();
      }
    }, 10000);
    this.liveCheck();
  }

  public liveCheck(): Subscription {
    return this.ping().pipe(first()).subscribe(
      data => {
        if (data) {
          this._isLiveSubject.next(true);
        } else {
          this._isLiveSubject.next(false);
        }
      },
      error => {
        console.log(error);
        this._isLiveSubject.next(false);
      }
    );
  }

  private fetchChannels(): void {
    if (this.isFetchingChannels) {
      return;
    }
    this.isFetchingChannels = true;
    this._isLoadingSubject.next(true);
    this.channels().pipe(first()).subscribe((value: ChannelResponse) => {
      if (value) {
        this._channelSubject.next(value);
        this.channelStore = value.data;
      }
      this.isFetchingChannels = false;
      setTimeout(() => this._isLoadingSubject.next(false), 250);
    });
  }

  getChannelById(id: number): Channel | null {
    if (this.channelStore.length === 0) {
      this.fetchChannels();
    }
    for (const channel of this.channelStore) {
      if (channel.id === id) {
        return channel;
      }
    }
    return null;
  }

  public statusResponse(): Observable<StatusResponse> {
    return this.get<StatusResponse>(this.apiEndpoint + '/status');
  }

  public channels(): Observable<ChannelResponse> {
    return this.get<ChannelResponse>(this.apiEndpoint + '/channels');
  }

  public updateChannelsOrder(channels: Channel[]): Observable<ChannelResponse> {
    return this.put<ChannelResponse>(this.apiEndpoint + '/channels', channels).pipe(
      tap((value: ChannelResponse) => {
        if (value) {
          this._channelSubject.next(value);
          this.channelStore = value.data;
        }
      })
    );
  }

  public program(from: Date, to: Date): Observable<ProgramResponse> {
    const url = `${this.apiEndpoint}/program`;
    const options = {
      params: {
        from: from.toISOString(),
        to: to.toISOString()
      }
    };
    return this.get<ProgramResponse>(url, options);
  }

  public entry(clickedEntryId: IdType): Observable<ProgramEntry> {
    const url = `${this.apiEndpoint}/program/entry/${encodeURIComponent(clickedEntryId.toString())}`;
    return this.get<ProgramEntry>(url);
  }

  public logEntries(page: number): Observable<LogEntryResponse> {
    const url = `${this.apiEndpoint}/log?page=${page}`;
    return this.get<LogEntryResponse>(url);
  }

  public recommendations(from: null | dayjs.Dayjs = null): Observable<Recommendation[]> {
    const url = `${this.apiEndpoint}/recommendations`;
    const options: any = {};
    if (from) {
      options.params = {
        from: from.toISOString()
      };
    }
    return this.get<Recommendation[]>(url, options);
  }

  public ping(): Observable<Pong> {
    return this.get<Pong>(this.apiEndpoint + '/ping');
  }

  public updateTvShows(){
    this.get<TvShow[]>(`${this.apiEndpoint}/tv-shows`)
      .pipe(first())
      .subscribe(value => this._tvShowSubject.next(value));
  }

  get tvShowSubject(): BehaviorSubject<TvShow[] | null> {
    return this._tvShowSubject;
  }

  get channelSubject(): BehaviorSubject<ChannelResponse | null> {
    return this._channelSubject;
  }

  get programSubject(): BehaviorSubject<ProgramResponse | null> {
    return this._programSubject;
  }

  get isLoadingSubject(): BehaviorSubject<boolean> {
    return this._isLoadingSubject;
  }

  get isInErrorsSubject(): BehaviorSubject<boolean> {
    return this._isInErrorsSubject;
  }

  get isLiveSubject(): BehaviorSubject<boolean | null> {
    return this._isLiveSubject;
  }

  get isWindowOpenedSubject(): BehaviorSubject<boolean> {
    return this._isWindowOpenedSubject;
  }

  get statusSubject(): BehaviorSubject<StatusResponse | null> {
    return this._statusSubject;
  }

  fetchProgramForDay(dateToFetch: Date): void {
    this._isLoadingSubject.next(true);

    this.fetchedDays.push(new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate()));

    const fromDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate(), 0, 0, 0, 0);
    const toDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() + 1, 23, 59, 59, 999);

    this.program(fromDate, toDate).pipe(first()).subscribe((value: ProgramResponse) => {
      if (!value) {
        return;
      }
      this._programSubject.next(value);
      setTimeout(() => this._isLoadingSubject.next(false), 1500);
    });
  }

  search(searchKey: string): Observable<ProgramEntry[]> {
    this._isLoadingSubject.next(true);
    return this.get<ProgramEntry[]>(this.apiEndpoint + '/search?query=' + encodeURIComponent(searchKey));
  }

  updateStatus(): void {
    this.statusResponse().pipe(first()).subscribe(statusResponse => {
      if (!statusResponse) {
        return;
      }
      this._statusSubject.next(statusResponse);
    });
  }

  /**
   * centralized http get with small error handling
   *
   * @param url
   * @param options
   * @private
   */
  private get<T>(url: string, options = {}): Observable<T> {
    if (!url.endsWith('/ping') && ((this._isInErrorsSubject.getValue() || this._isLiveSubject.getValue() === false))) {
      console.log(`api in errors or not live. Skipping request to url ${url}.`);
      return EMPTY;
    }
    const inErrAlready = this._isInErrorsSubject.getValue();
    return this.http.get<T>(url, options).pipe(
      timeout(environment.apiRequestTimeoutInSecs * 1000),
      tap(
        _ => {
          if (inErrAlready) {
            this._isInErrorsSubject.next(false);
          }
        }
      ),
      catchError(err => {
        if (err.name === 'TimeoutError') {
          console.log('request timeout reached!', err);
        }
        if (!inErrAlready) {
          this._isInErrorsSubject.next(true);
        }
        if (url.endsWith('/ping')) {
          this._isLiveSubject.next(false);
        }
        console.error('http GET call err', url, err);
        return EMPTY;
      })
    );
  }

  /**
   * centralized http put with small error handling
   *
   * @param url
   * @param body
   * @param options
   * @private
   */
  private put<T>(url: string, body: any, options = {}): Observable<T> {
    if (this._isInErrorsSubject.getValue() || this._isLiveSubject.getValue() === false) {
      console.log(`api in errors or not live. Skipping request to url ${url}.`);
      return EMPTY;
    }
    const inErrAlready = this._isInErrorsSubject.getValue();
    return this.http.put<T>(url, body, options).pipe(
      timeout(environment.apiRequestTimeoutInSecs * 1000),
      tap(
        _ => {
          if (inErrAlready) {
            this._isInErrorsSubject.next(false);
          }
        }
      ),
      catchError(err => {
        if (err.name === 'TimeoutError') {
          console.log('request timeout reached!', err);
        }
        if (!inErrAlready) {
          this._isInErrorsSubject.next(true);
        }
        console.error('http PUT call err', url, err);
        return EMPTY;
      })
    );
  }
}
