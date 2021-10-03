import { Pipe, PipeTransform } from '@angular/core';
import moment, {MomentInput} from 'moment-timezone';
import {environment} from '../../environments/environment';

@Pipe({
  name: 'momentDate'
})
export class MomentDatePipe implements PipeTransform {

  public readonly FULL_DATE_TIME_HUMAN = 'D.MM.YYYY - HH:mm:ss';
  public readonly MEDIUM_TIME = 'HH:mm';
  public readonly MEDIUM_DATE_TIME = 'D.MM.YYYY - HH:mm';

  transform(value: MomentInput, ...args: string[]): string {
    if (args.length === 0) {
      console.error('empty value for momentDate pipe received!');
      return 'empty';
    }
    let dateFormat = this.FULL_DATE_TIME_HUMAN;
    switch (args[0]) {
      case 'full': break;
      case 'time-medium':
        dateFormat = this.MEDIUM_TIME; break;
      case 'full-medium':
        dateFormat = this.MEDIUM_DATE_TIME; break;
    }
    return moment(value).tz(environment.timezone).format(dateFormat);
  }

}
