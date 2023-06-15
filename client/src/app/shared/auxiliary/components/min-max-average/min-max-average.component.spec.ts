import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MinMaxAverageComponent } from './min-max-average.component';

describe('MinMaxAverageComponent', () => {
  let component: MinMaxAverageComponent;
  let fixture: ComponentFixture<MinMaxAverageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ MinMaxAverageComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(MinMaxAverageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
