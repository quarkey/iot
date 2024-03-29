import { Component, Input, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import {
  Controller,
  normalswitch,
  thresholdswitch,
  timeswitch,
  webcamstreamtimelapse,
} from 'src/app/models/controllers';
import { Sensordata } from 'src/app/models/dataset';
import { ControllersService } from 'src/app/services/controllers.service';
import { DialogsService } from 'src/app/services/dialogs.service';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-controller-details',
  templateUrl: './controller-details.component.html',
  styleUrls: ['./controller-details.component.scss'],
})
export class ControllerDetailsComponent implements OnInit {
  @Input() citem: Controller;
  form: FormGroup;
  showReloadbutton = false;
  categories: string[] = ['switch', 'thresholdswitch', 'timeswitch', 'timeswitchrepeat', 'webcamstreamtimelapse'];

  constructor(
    private formBuilder: FormBuilder,
    private controllerService: ControllersService,
    private dialogService: DialogsService,
    private router: Router
  ) {}
  defaultValue = { hour: 13, minute: 30 };

  timeChangeHandler(event: Event) {
    console.log(event);
  }

  invalidInputHandler() {
    // some error handling
  }
  ngOnInit(): void {
    this.form = this.controllerService.addInitialForm(this.citem);
    this.citem.items.forEach((item: any) => {
      switch (this.citem.category) {
        case 'thresholdswitch':
          this.items.push(this.controllerService.addThresholdSwitchForm(item as thresholdswitch));
          break;
        case 'timeswitch':
          this.items.push(this.controllerService.addTimeSwitchForm(item as timeswitch));
          break;
        case 'timeswitchrepeat':
          this.items.push(this.controllerService.addTimeSwitchRepeatForm(item as timeswitch));
          break;
        case 'switch':
          this.items.push(this.controllerService.addSwitchForm(item as normalswitch));
          break;
        case 'webcamstreamtimelapse':
          this.items.push(this.controllerService.addWebcamStreamTimelapse(item as webcamstreamtimelapse));
          break;
      }
    });
    this.form.controls.category.valueChanges.subscribe((value) => {
      this.showReloadbutton = true;
      this.confirmReload();
    });
  }
  get items() {
    return this.form.get('items') as FormArray;
  }
  updateController() {
    var obj = {
      ...this.form.value,
      id: this.citem.id,
    };

    this.controllerService.UpdateControllerByID(obj).subscribe((res) => {
      if (res) {
        this.form.markAsPristine();
      }
    });
  }
  confirmReload() {
    this.dialogService
      .openConfirmationDialog(
        'Save and reload page?',
        `Warning! To initiate new form values the page must be reloaded. 
        Do you want to save values and reload page?`,
        'CONFIRM',
        'CANCEL',
        true
      )
      .subscribe((res) => {
        if (res) {
          this.showReloadbutton = false;
          this.updateController(); // save form
          window.location.reload();
        }
      });
  }
  deleteController() {
    this.dialogService
      .openConfirmationDialog(
        'Delete controller?',
        `Are you sure you want to permanently delete controller?`,
        'CONFIRM',
        'CANCEL',
        true
      )
      .subscribe((res) => {
        if (res) {
          this.showReloadbutton = false;
          // do delete
          this.controllerService.DeleteControllerByID(this.citem.id).subscribe((res) => {
            if (res) {
              this.router.navigate([`/controllers`]);
            }
          });
        }
      });
  }
  resetControllerItemValues() {
    this.dialogService
      .openConfirmationDialog(
        'Reset item fields?',
        `Are you sure you want to reset item values?`,
        'CONFIRM',
        'CANCEL',
        false
      )
      .subscribe((res) => {
        if (res) {
          this.showReloadbutton = false;
          // do reset
          this.controllerService
            .ResetControllerSwitchValueEndpoint(this.citem.id, this.citem.category)
            .subscribe((res) => {
              if (res) {
                window.location.reload();
              }
            });
        }
      });
  }
}
