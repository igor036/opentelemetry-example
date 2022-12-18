import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { AddressModule } from './address/address.module';
import { AppDatasourceTest } from './db/config';

@Module({
  imports: [
    AddressModule,
    TypeOrmModule.forRoot(AppDatasourceTest),
  ],
  controllers: [],
  providers: [],
})
export class AppTestModule {}
