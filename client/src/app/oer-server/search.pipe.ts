/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2024 emschu[aet]mailbox.org
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
import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'search'
})
export class SearchPipe implements PipeTransform {

  transform(value: string | undefined, ...args: string[]): string {
    if (value && args.length > 0) {
      const split = args[0].split(',');
      for (const singleKeyword of split) {
        value = value?.replace(new RegExp('(' + singleKeyword + ')', 'gi'), (a: string, b: string) => {
          return '<span class="keyword-search">' + b + '</span>';
        });
      }
      return value;
    }
    return '';
  }
}
