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
import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {Recommendation} from '../entities';
import {ApiService} from '../api.service';
import {AbstractReadMoreComponent} from '../AbstractReadMoreComponent';
import {Subscription} from 'rxjs';
import {first} from 'rxjs/operators';
import dayjs from 'dayjs';

@Component({
    selector: 'app-recommendation',
    templateUrl: './recommendation.component.html',
    styleUrls: ['./recommendation.component.scss'],
    standalone: false
})
export class RecommendationComponent extends AbstractReadMoreComponent implements OnInit, OnDestroy {

  recommendations: Recommendation[] | null = null;

  private recommendationSubscription: Subscription | null = null;

  constructor(public apiService: ApiService) {
    super();
  }

  ngOnInit(): void {
    this.loadRecommendations();
  }

  ngOnDestroy(): void {
    this.recommendationSubscription?.unsubscribe();
  }

  private loadRecommendations(): void {
    this.apiService.isLoadingSubject.next(true);
    this.fetchRecommendations('now');
  }

  @HostListener('click', ['$event'])
  onClick(e: any): void {
    if (e?.target?.classList.contains('read-more')) {
      this.onReadMore(e);
    }
  }

  fetchRecommendations(timeExpression: string): void {
    let from: dayjs.Dayjs = dayjs();
    let isNow = false;

    switch (timeExpression) {
      case 'now':
        isNow = true;
        break;
      case 'tomorrow':
        from = from.add(1, 'day');
        from = from.hour(8);
        break;
      case 'dayAfterTomorrow':
        from = from.add(2, 'day');
        from = from.hour(0);
        break;
      case '20':
        from = from.hour(20);
        break;
      case '22':
        from = from.hour(22);
        break;
      case '0':
        from = from.add(1, 'day');
        from = from.hour(0);
        break;
    }
    if (!isNow) {
      from = from.minute(0);
      from = from.second(0);
      from = from.millisecond(0);
    }
    this.apiService.isLoadingSubject.next(true);
    this.apiService.recommendations(from).pipe(first()).subscribe(value => {
      this.recommendations = value;
      setTimeout(() => {
        this.apiService.isLoadingSubject.next(false);
      }, 250);
    });
  }
}
