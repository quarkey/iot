import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SelectDeviceComponent } from './select-device.component';

describe('SelectDeviceComponent', () => {
  let component: SelectDeviceComponent;
  let fixture: ComponentFixture<SelectDeviceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SelectDeviceComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SelectDeviceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
