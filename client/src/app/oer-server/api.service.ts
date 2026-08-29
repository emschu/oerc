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
import {Injectable, signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {EMPTY, Observable, Subscription} from 'rxjs';
import {
  Channel,
  ChannelResponse,
  LogEntryResponse,
  Pong,
  ProgramEntry,
  ProgramResponse,
  Recommendation,
  StatusResponse,
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

  private _isLiveSubject = signal<boolean | null>(null);
  public isLive = this._isLiveSubject.asReadonly();
  private _channelSubject = signal<ChannelResponse | null>(null);
  public channelsSignal = this._channelSubject.asReadonly();
  private _programSubject = signal<ProgramResponse | null>(null);
  public programSignal = this._programSubject.asReadonly();
  private _tvShowSubject = signal<TvShow[] | null>(null);
  public tvShows = this._tvShowSubject.asReadonly();
  private _isLoadingSubject = signal<boolean>(true);
  public isLoading = this._isLoadingSubject.asReadonly();
  private _isInErrorsSubject = signal<boolean>(false);
  public isInErrors = this._isInErrorsSubject.asReadonly();
  private _isWindowOpenedSubject = signal<boolean>(true);
  public isWindowOpened = this._isWindowOpenedSubject.asReadonly();
  private _statusSubject = signal<StatusResponse | null>(null);
  public status = this._statusSubject.asReadonly();

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
      if (this._isWindowOpenedSubject()) {
        this.liveCheck();
      }
    }, 10000);
    this.liveCheck();
  }

  public liveCheck(): Subscription {
    return this.ping().pipe(first()).subscribe(
      data => {
        if (data) {
          this._isLiveSubject.set(true);
        } else {
          this._isLiveSubject.set(false);
        }
      },
      error => {
        console.log(error);
        this._isLiveSubject.set(false);
      }
    );
  }

  private fetchChannels(): void {
    if (this.isFetchingChannels) {
      return;
    }
    this.isFetchingChannels = true;
    this._isLoadingSubject.set(true);
    this.channels().pipe(first()).subscribe((value: ChannelResponse) => {
      if (value) {
        this._channelSubject.set(value);
        this.channelStore = value.data;
      }
      this.isFetchingChannels = false;
      setTimeout(() => this._isLoadingSubject.set(false), 250);
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
          this._channelSubject.set(value);
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

  public updateTvShows() {
    this.get<TvShow[]>(`${this.apiEndpoint}/tv-shows`)
      .pipe(first())
      .subscribe(value => this._tvShowSubject.set(value));
  }

  get tvShowSubject(): WritableSignal<TvShow[] | null> {
    return this._tvShowSubject;
  }

  get channelSubject(): WritableSignal<ChannelResponse | null> {
    return this._channelSubject;
  }

  get programSubject(): WritableSignal<ProgramResponse | null> {
    return this._programSubject;
  }

  get isLoadingSubject(): WritableSignal<boolean> {
    return this._isLoadingSubject;
  }

  get isInErrorsSubject(): WritableSignal<boolean> {
    return this._isInErrorsSubject;
  }

  get isLiveSubject(): WritableSignal<boolean | null> {
    return this._isLiveSubject;
  }

  get isWindowOpenedSubject(): WritableSignal<boolean> {
    return this._isWindowOpenedSubject;
  }

  get statusSubject(): WritableSignal<StatusResponse | null> {
    return this._statusSubject;
  }

  fetchProgramForDay(dateToFetch: Date): void {
    this._isLoadingSubject.set(true);

    this.fetchedDays.push(new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate()));

    const fromDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate(), 0, 0, 0, 0);
    const toDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() + 1, 23, 59, 59, 999);

    this.program(fromDate, toDate).pipe(first()).subscribe((value: ProgramResponse) => {
      if (!value) {
        return;
      }
      this._programSubject.set(value);
      setTimeout(() => this._isLoadingSubject.set(false), 1500);
    });
  }

  search(searchKey: string): Observable<ProgramEntry[]> {
    this._isLoadingSubject.set(true);
    return this.get<ProgramEntry[]>(this.apiEndpoint + '/search?query=' + encodeURIComponent(searchKey));
  }

  updateStatus(): void {
    this.statusResponse().pipe(first()).subscribe(statusResponse => {
      if (!statusResponse) {
        return;
      }
      this._statusSubject.set(statusResponse);
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
    if (!url.endsWith('/ping') && ((this._isInErrorsSubject() || this._isLiveSubject() === false))) {
      console.log(`api in errors or not live. Skipping request to url ${url}.`);
      return EMPTY;
    }
    const inErrAlready = this._isInErrorsSubject();
    return this.http.get<T>(url, options).pipe(
      timeout(environment.apiRequestTimeoutInSecs * 1000),
      tap(
        _ => {
          if (inErrAlready) {
            this._isInErrorsSubject.set(false);
          }
        }
      ),
      catchError(err => {
        if (err.name === 'TimeoutError') {
          console.log('request timeout reached!', err);
        }
        if (!inErrAlready) {
          this._isInErrorsSubject.set(true);
        }
        if (url.endsWith('/ping')) {
          this._isLiveSubject.set(false);
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
    if (this._isInErrorsSubject() || this._isLiveSubject() === false) {
      console.log(`api in errors or not live. Skipping request to url ${url}.`);
      return EMPTY;
    }
    const inErrAlready = this._isInErrorsSubject();
    return this.http.put<T>(url, body, options).pipe(
      timeout(environment.apiRequestTimeoutInSecs * 1000),
      tap(
        _ => {
          if (inErrAlready) {
            this._isInErrorsSubject.set(false);
          }
        }
      ),
      catchError(err => {
        if (err.name === 'TimeoutError') {
          console.log('request timeout reached!', err);
        }
        if (!inErrAlready) {
          this._isInErrorsSubject.set(true);
        }
        console.error('http PUT call err', url, err);
        return EMPTY;
      })
    );
  }
}
