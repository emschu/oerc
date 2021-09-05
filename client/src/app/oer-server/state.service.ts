import { Injectable } from '@angular/core';
import {environment} from '../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class StateService {

  private readonly KEY_SHOW_DEPRECATED = 'show_deprecated_entries';

  constructor() { }

  getShowDeprecatedEntries(): boolean {
    if (localStorage.getItem(this.KEY_SHOW_DEPRECATED) === null) {
      this.setShowDeprecatedEntries(environment.defaultSettingShowDeprecatedEntries);
    }
    return JSON.parse(localStorage.getItem(this.KEY_SHOW_DEPRECATED) as string);
  }

  setShowDeprecatedEntries(value: boolean): void {
    localStorage.setItem(this.KEY_SHOW_DEPRECATED, JSON.stringify(value));
  }
}
