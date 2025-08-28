import { Routes } from '@angular/router';
import { LoginComponent } from './pages/login.component';
import { RegisterComponent } from './pages/register.component';
import { ConcertsComponent } from './pages/concerts.component';

export const routes: Routes = [
  { path: '', redirectTo: 'concerts', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'concerts', component: ConcertsComponent },
  { path: '**', redirectTo: 'concerts' },
];
