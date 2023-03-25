import { Component, Input } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { Device } from 'src/app/models/device';

@Component({
  selector: 'app-select-device',
  templateUrl: './select-device.component.html',
  styleUrls: ['./select-device.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      multi: true,
      useExisting: SelectDeviceComponent,
    },
  ],
})
export class SelectDeviceComponent implements ControlValueAccessor {
  @Input() deviceList: Device[];
  sensor_id: any;
  onChange = (sensor_id) => {};
  onTouched = () => {};
  touched = false;
  disabled = false;
  change: any;

  handleChange() {
    this.onChange(this.change);
  }

  writeValue(sensor_id: any): void {
    this.sensor_id = sensor_id;
  }
  registerOnChange(onChange: any): void {
    this.onChange = onChange;
  }
  registerOnTouched(onTouched: any): void {
    this.onTouched = onTouched;
  }
  setDisabledState?(disabled: boolean): void {
    this.disabled = disabled;
  }
  markAsTouched() {
    if (!this.touched) {
      this.onTouched();
      this.touched = true;
    }
  }
}
