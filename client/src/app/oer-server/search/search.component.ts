import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {ApiService} from '../api.service';
import {ProgramEntry} from '../entities';
import {AbstractReadMoreComponent} from '../AbstractReadMoreComponent';
import {SearchService} from './search.service';
import {Observable, Subscription} from 'rxjs';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss']
})
export class SearchComponent extends AbstractReadMoreComponent implements OnInit, OnDestroy {
  programEntryList: ProgramEntry[] = [];
  searchString = '';
  isSearchInProgressIndicator = false;
  isInErrors = false;

  private searchSubscription: Subscription | null = null;

  constructor(private activeRoute: ActivatedRoute,
              private searchService: SearchService,
              public apiService: ApiService) {
    super();
  }

  ngOnInit(): void {
    this.initSearch();
  }

  private initSearch(): void {
    this.activeRoute.queryParamMap.subscribe((value) => {
      const searchKey = value.get('query') ?? '';
      this.searchString = searchKey;
      this.searchService.lastSearchStringSubject.next(searchKey);
      this.isSearchInProgressIndicator = true;
      this.apiService.isLoadingSubject.next(true);
      this.searchSubscription = this.apiService.search(searchKey).subscribe(value1 => {
          this.isInErrors = false;
          if (value1) {
            this.programEntryList = value1;
          } else {
            this.programEntryList = [];
          }
          this.isSearchInProgressIndicator = false;
          setTimeout(() => {
            this.apiService.isLoadingSubject.next(false);
          }, 250);
        },
        err => {
          this.isInErrors = true;
          console.error(err);
          this.programEntryList = [];
          this.isSearchInProgressIndicator = false;
          setTimeout(() => {
            this.apiService.isLoadingSubject.next(false);
          }, 250);
          return new Observable();
        });
    });
  }

  ngOnDestroy(): void {
    this.searchSubscription?.unsubscribe();
  }

  @HostListener('click', ['$event'])
  onClick(e: any): void {
    if (e?.target?.classList.contains('read-more')) {
      this.onReadMore(e);
    }
  }
}
