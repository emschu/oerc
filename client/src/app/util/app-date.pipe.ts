/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2025 emschu[aet]mailbox.org
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
import { Pipe, PipeTransform } from '@angular/core';
import {environment} from '../../environments/environment';
import dayjs, {Dayjs} from 'dayjs';

@Pipe({
    name: 'appDate',
    standalone: false
})
export class AppDatePipe implements PipeTransform {

  public readonly FULL_DATE_TIME_HUMAN = 'D.MM.YYYY - HH:mm:ss';
  public readonly MEDIUM_DATE = 'D.MM.YYYY';
  public readonly MEDIUM_TIME = 'HH:mm';
  public readonly MEDIUM_DATE_TIME = 'D.MM.YYYY - HH:mm';

  transform(value: dayjs.ConfigType, ...args: string[]): string {
    if (args.length === 0) {
      console.error('empty value for appDate pipe received!');
      return 'empty';
    }
    let dateFormat = this.FULL_DATE_TIME_HUMAN;
    switch (args[0]) {
      case 'full': break;
      case 'time-medium':
        dateFormat = this.MEDIUM_TIME; break;
      case 'date-medium':
        dateFormat = this.MEDIUM_DATE; break;
      case 'full-medium':
        dateFormat = this.MEDIUM_DATE_TIME; break;
    }
    return dayjs(value).locale(environment.locale).format(dateFormat);
  }
}
