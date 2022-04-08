import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ToolbarService {
  pageTitle = 'default';
  constructor() {}

  setTitle(title: string) {
    this.pageTitle = title;
  }
  getTitle() {
    return this.pageTitle;
  }
}
