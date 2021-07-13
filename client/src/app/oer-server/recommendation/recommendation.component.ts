import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {Recommendation} from '../entities';
import {ApiService} from '../api.service';
import {AbstractReadMoreComponent} from '../AbstractReadMoreComponent';
import {environment} from '../../../environments/environment';
import moment from 'moment-timezone';
import {Subscription} from 'rxjs';

@Component({
  selector: 'app-recommendation',
  templateUrl: './recommendation.component.html',
  styleUrls: ['./recommendation.component.scss']
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

  fetchRecommendations(now: string): void {
    const currentOffset = moment(new Date()).tz(environment.timezone).utcOffset();
    let from: moment.Moment = moment(new Date()).tz(environment.timezone).utc(true);
    let isNow = false;

    switch (now) {
      case 'now':
        isNow = true;
        break;
      case 'tomorrow':
        from = from.add(1, 'day');
        from = from.hours(8);
        break;
      case 'dayAfterTomorrow':
        from = from.add(2, 'day');
        from = from.hours(0);
        break;
      case '20':
        from = from.hours(20);
        break;
      case '22':
        from = from.hours(22);
        break;
      case '0':
        from = from.add(1, 'day');
        from = from.hours(0);
        break;
    }
    if (!isNow) {
      from = from.minute(0);
      from = from.seconds(0);
      from = from.millisecond(0);
    }
    this.apiService.isLoadingSubject.next(true);
    this.apiService.recommendations(from).subscribe(value => {
      this.recommendations = value;
      setTimeout(() => {
        this.apiService.isLoadingSubject.next(false);
      }, 250);
    });
  }
}
