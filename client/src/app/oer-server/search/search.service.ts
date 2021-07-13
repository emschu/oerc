import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SearchService {
  private _lastSearchStringSubject = new BehaviorSubject<string>('');

  constructor() {
  }

  get lastSearchStringSubject(): BehaviorSubject<string> {
    return this._lastSearchStringSubject;
  }
}
