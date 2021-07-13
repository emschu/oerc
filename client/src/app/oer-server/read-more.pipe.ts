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
import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'readMore'
})
export class ReadMorePipe implements PipeTransform {

  private maxCharsVisible = 750;

  transform(value: string | undefined, ...args: any[]): string | undefined {
    if (value && value.toString().length > this.maxCharsVisible) {
      let output = '<div class="d-inline-block">' + value.substr(0, this.maxCharsVisible) +
        '<a class="read-more btn btn-link btn-sm">... More</a>';
      output += '<div class="d-inline d-hide">' + value.substr(this.maxCharsVisible) + '</div>';
      output += ' <a class="read-less btn btn-link btn-sm d-hide">Show Less</a>';
      output += '</div>';
      return output;
    }
    return value;
  }
}
