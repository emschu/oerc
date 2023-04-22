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
