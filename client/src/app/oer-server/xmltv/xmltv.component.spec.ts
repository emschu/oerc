import { ComponentFixture, TestBed } from '@angular/core/testing';

import { XmltvComponent } from './xmltv.component';

describe('XmltvComponent', () => {
  let component: XmltvComponent;
  let fixture: ComponentFixture<XmltvComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [XmltvComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(XmltvComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
