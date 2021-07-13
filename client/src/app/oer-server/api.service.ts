import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BehaviorSubject, interval, Observable, Subscription} from 'rxjs';
import {Channel, ChannelResponse, LogEntryResponse, Pong, ProgramEntry, ProgramResponse, Recommendation, StatusResponse} from './entities';
import {IdType} from 'vis-timeline';
import {catchError, tap, timeout} from 'rxjs/operators';
import {environment} from '../../environments/environment';
import {Moment} from 'moment-timezone';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiEndpoint = `${environment.serverEndpoint}`;

  private _isLiveSubject = new BehaviorSubject<boolean | null>(null);
  private _channelSubjectVar = new BehaviorSubject<ChannelResponse | null>(null);
  private _programSubject = new BehaviorSubject<ProgramResponse | null>(null);
  private _isLoadingSubject = new BehaviorSubject<boolean>(true);
  private _isInErrorsSubject = new BehaviorSubject<boolean>(false);
  private _isWindowOpenedSubject = new BehaviorSubject<boolean>(true);

  private channelStore: Channel[] = [];
  private fetchedDays: Date[] = [];
  private isFetchingChannels = false;

  constructor(public http: HttpClient) {
    // initial fetch of channels
    this.liveCheck();

    this.fetchChannels();

    interval(10000).subscribe(() => {
      if (this._isWindowOpenedSubject.getValue()) {
        this.liveCheck.bind(this);
      }
    });
  }

  public liveCheck(): Subscription {
    return this.ping().subscribe(
      data => {
        if (!data) {
          this._isLiveSubject.next(false);
        } else {
          this._isLiveSubject.next(true);
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
    this.isLoadingSubject.next(true);
    this.channels().subscribe((value: ChannelResponse) => {
      if (value) {
        this._channelSubjectVar.next(value);
        this.channelStore = value.data;
      }
      this.isFetchingChannels = false;
      setTimeout(() => this.isLoadingSubject.next(false), 250);
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

  public dailyProgram(): Observable<ProgramResponse> {
    return this.get<ProgramResponse>(this.apiEndpoint + '/program/daily');
  }

  public program(from: Date, to: Date): Observable<ProgramResponse> {
    return this.get<ProgramResponse>(this.apiEndpoint + '/program?from=' + from.toISOString() + '&to=' + to.toISOString());
  }

  public entry(clickedEntryId: IdType): Observable<ProgramEntry> {
    return this.get<ProgramEntry>(this.apiEndpoint + '/program/entry/' + clickedEntryId);
  }

  public logEntries(offset: number = 0, limit: number = 500): Observable<LogEntryResponse> {
    return this.get<LogEntryResponse>(this.apiEndpoint + '/log');
  }

  public recommendations(from: null | Moment = null): Observable<Recommendation[]> {
    let queryParams = '';
    if (from) {
      from = from.tz(environment.timezone).utc(false);
      queryParams += 'from=' + encodeURIComponent(from.toISOString());
    }
    if (queryParams.length > 0) {
      return this.get<Recommendation[]>(this.apiEndpoint + '/recommendations?' + queryParams);
    }
    return this.get<Recommendation[]>(this.apiEndpoint + '/recommendations');
  }

  public ping(): Observable<Pong> {
    return this.get<Pong>(this.apiEndpoint + '/ping');
  }

  get channelSubjectVar(): BehaviorSubject<ChannelResponse | null> {
    return this._channelSubjectVar;
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

  fetchProgramForDay(dateToFetch: Date): void {
    this.isLoadingSubject.next(true);

    this.fetchedDays.push(new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate()));

    const fromDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() - 1, 0, 0, 0, 0);
    const toDate: Date = new Date(dateToFetch.getFullYear(), dateToFetch.getMonth(), dateToFetch.getDate() + 1, 23, 59, 59, 999);

    this.program(fromDate, toDate).subscribe((value: ProgramResponse) => {
      if (!value) {
        return;
      }
      this._programSubject.next(value);
      setTimeout(() => this.isLoadingSubject.next(false), 1500);
    });
  }

  checkIfDayIsFetched(dayToCheck: Date): boolean {
    let isFetched = false;
    this.fetchedDays.forEach(each => {
      if (each.getDate() === dayToCheck.getDate() && each.getMonth() === dayToCheck.getMonth()
        && each.getFullYear() === dayToCheck.getFullYear()) {
        isFetched = true;
      }
    });
    return isFetched;
  }

  search(searchKey: string): Observable<ProgramEntry[]> {
    this.isLoadingSubject.next(true);
    return this.get<ProgramEntry[]>(this.apiEndpoint + '/search?query=' + encodeURIComponent(searchKey));
  }

  /**
   * centralized http get with small error handling
   *
   * @param url
   * @param options
   * @private
   */
  private get<T>(url: string, options = {}): Observable<T> {
    if (!url.endsWith('/ping') && (this.isInErrorsSubject.getValue() || this.isLiveSubject.getValue() === false)) {
      console.log(`api in errors or not live. Skipping request to url ${url}.`);
      return new Observable<T>();
    }
    const inErrAlready = this._isInErrorsSubject.getValue();
    return this.http.get<T>(url).pipe(
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
          throw err;
        }
        if (!inErrAlready) {
          this._isInErrorsSubject.next(true);
        }
        if (url.endsWith('/ping')) {
          this._isLiveSubject.next(false);
        }
        console.error('http GET call err', url, err);
        return new Observable<T>();
      })
    );
  }
}
