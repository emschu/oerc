import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {ApiService} from '../oer-server/api.service';
import {SearchService} from '../oer-server/search/search.service';
import {Router} from '@angular/router';

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss']
})
export class NavComponent implements OnInit, OnDestroy {
  currentSearchPhrase = '';
  isLoading = false;

  private searchPhraseSubscription: Subscription|null = null;
  private isLoadingSubscription: Subscription|null = null;

  constructor(public apiService: ApiService,
              public searchService: SearchService,
              private router: Router) { }

  ngOnInit(): void {
    this.searchPhraseSubscription = this.searchService.lastSearchStringSubject.subscribe(value => {
      if (value !== this.currentSearchPhrase) {
        this.currentSearchPhrase = value;
      }
    });
    this.isLoadingSubscription = this.apiService.isLoadingSubject.subscribe(value => {
      this.isLoading = value;
    });
  }

  searchFor(): void {
    const elementById = document.getElementById('search_text');
    if (elementById instanceof HTMLInputElement) {
      const searchWord = elementById?.value;
      if (searchWord.length > 2) {
        this.router.navigate(['/search'], {queryParams: {query: searchWord}});
      }
    }
  }

  openSection(): void {
    this.searchService.lastSearchStringSubject.next('');
  }

  ngOnDestroy(): void {
    this.searchPhraseSubscription?.unsubscribe();
    this.isLoadingSubscription?.unsubscribe();
  }
}
