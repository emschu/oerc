/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2025 emschu[aet]mailbox.org
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
import { waitForAsync, ComponentFixture, TestBed } from '@angular/core/testing';
import { ApiService } from '../api.service';
import { of, BehaviorSubject } from 'rxjs';

import { StatusComponent } from './status.component';

describe('StatusComponent', () => {
  let component: StatusComponent;
  let fixture: ComponentFixture<StatusComponent>;
  let apiServiceMock: any;

  beforeEach(waitForAsync(() => {
    apiServiceMock = {
      ping: jasmine.createSpy('ping').and.returnValue(of({})),
      updateStatus: jasmine.createSpy('updateStatus'),
      isLiveSubject: new BehaviorSubject<boolean | null>(null),
      isInErrorsSubject: new BehaviorSubject<boolean>(false),
      statusSubject: new BehaviorSubject<any>(null),
      isWindowOpenedSubject: new BehaviorSubject<boolean>(true)
    };

    TestBed.configureTestingModule({
      declarations: [ StatusComponent ],
      providers: [
        { provide: ApiService, useValue: apiServiceMock }
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StatusComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
