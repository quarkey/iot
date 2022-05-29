import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TwoColumnDivComponent } from './two-column-div.component';

describe('TwoColumnDivComponent', () => {
  let component: TwoColumnDivComponent;
  let fixture: ComponentFixture<TwoColumnDivComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TwoColumnDivComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TwoColumnDivComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
