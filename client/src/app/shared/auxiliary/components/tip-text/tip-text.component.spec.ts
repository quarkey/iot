import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TipTextComponent } from './tip-text.component';

describe('TipTextComponent', () => {
  let component: TipTextComponent;
  let fixture: ComponentFixture<TipTextComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ TipTextComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(TipTextComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
