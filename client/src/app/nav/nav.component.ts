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

  private searchPhraseSubscription: Subscription | null = null;
  private isLoadingSubscription: Subscription | null = null;
  private searchTextElement: HTMLElement | null = null;

  constructor(public apiService: ApiService,
              public searchService: SearchService,
              private router: Router) {
  }

  ngOnInit(): void {
    this.searchTextElement = document.getElementById('search_text');

    this.searchPhraseSubscription = this.searchService.lastSearchStringSubject.subscribe(value => {
      if (value !== this.currentSearchPhrase) {
        this.currentSearchPhrase = value;
      }
    });
    this.isLoadingSubscription = this.apiService.isLoadingSubject.subscribe(value => {
      this.isLoading = value;
    });
    this.searchTextElement?.focus();
  }

  searchFor(): void {
    if (this.searchTextElement instanceof HTMLInputElement) {
      const searchWord = this.searchTextElement?.value;
      if (searchWord.length > 2) {
        this.router.navigate(['/search'], {queryParams: {query: searchWord}, skipLocationChange: false, onSameUrlNavigation: 'reload'});
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
