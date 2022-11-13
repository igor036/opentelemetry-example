import { Body, Controller, Post } from '@nestjs/common';
import { Address } from './address.model';
import { AddressService } from './address.service';

@Controller('address')
export class AddressController {

    constructor(
        private addressService: AddressService
    ) {}

    @Post()
    async create(@Body() address: Address): Promise<Address> {
        return await this.addressService.save(address);
    }

}