/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
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
import {ChangeDetectionStrategy, Component, effect, inject, OnInit} from '@angular/core';
import {ApiService} from '../oer-server/api.service';
import {SearchService} from '../oer-server/search/search.service';
import {Router, RouterLink} from '@angular/router';
import {NgClass} from '@angular/common';

@Component({
  selector: 'app-nav',
  templateUrl: './nav.component.html',
  styleUrls: ['./nav.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true,
  imports: [
    RouterLink,
    NgClass
  ]
})
export class NavComponent implements OnInit {
  currentSearchPhrase = '';

  public apiService = inject(ApiService);
  public searchService = inject(SearchService);
  private router = inject(Router);

  private searchTextElement: HTMLElement | null = null;

  constructor() {
    effect(() => {
      this.currentSearchPhrase = this.searchService.lastSearchString();
    });
  }

  ngOnInit(): void {
    this.searchTextElement = document.getElementById('search_text');
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
    this.searchService.lastSearchStringSubject.set('');
  }
}
