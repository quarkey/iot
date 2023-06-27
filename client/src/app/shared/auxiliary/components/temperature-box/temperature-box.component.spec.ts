import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TemperatureBoxComponent } from './temperature-box.component';

describe('TemperatureBoxComponent', () => {
  let component: TemperatureBoxComponent;
  let fixture: ComponentFixture<TemperatureBoxComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TemperatureBoxComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TemperatureBoxComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
