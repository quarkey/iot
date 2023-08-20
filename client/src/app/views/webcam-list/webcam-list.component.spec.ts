import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WebcamListComponent } from './webcam-list.component';

describe('WebcamListComponent', () => {
  let component: WebcamListComponent;
  let fixture: ComponentFixture<WebcamListComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [WebcamListComponent]
    });
    fixture = TestBed.createComponent(WebcamListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
