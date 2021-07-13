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
