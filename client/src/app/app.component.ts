import {ApiService} from './oer-server/api.service';
import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {OnPageHidden, OnPageVisible} from 'angular-page-visibility';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  isLive = false;
  private isLiveSubscription: Subscription | null = null;
  private inited = false;

  constructor(public apiService: ApiService) {
  }

  @OnPageHidden()
  onPageHidden(): void {
    console.log('on page hidden');
    this.apiService.isWindowOpenedSubject.next(false);
  }

  @OnPageVisible()
  onPageVisible(): void {
    console.log('on page visible');
    this.apiService.isWindowOpenedSubject.next(true);
  }

  ngOnInit(): void {
    this.isLiveSubscription = this.apiService.isLiveSubject.subscribe(value => {
      if (this.inited && !this.isLive && value !== null) {
        window.location.reload();
      }
      if (!this.inited && value !== null) {
        this.inited = true;
      }
      if (value === true || value === false) {
        this.isLive = value;
      }
    });
  }

  ngOnDestroy(): void {
    this.isLiveSubscription?.unsubscribe();
  }
}
