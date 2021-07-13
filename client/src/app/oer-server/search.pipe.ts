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
