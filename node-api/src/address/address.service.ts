import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Address } from './address.model';

@Injectable()
export class AddressService {

    constructor(
        @InjectRepository(Address)
        private addressRepository: Repository<Address>,
    ) {}

    save(address: Address): Promise<Address> {
        return this.addressRepository.save(address);
    }
}