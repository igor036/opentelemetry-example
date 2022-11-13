import { IsNotEmpty, MaxLength, MinLength } from 'class-validator';
import { Entity, Column, PrimaryGeneratedColumn } from 'typeorm';

@Entity()
export class Address {
    
    @PrimaryGeneratedColumn()
    id: number;

    @IsNotEmpty()
    @MinLength(8)
    @MaxLength(8)
    @Column({unique: true})
    zipCode: string;

    @Column()
    @IsNotEmpty()
    address: string;

    @Column()
    @IsNotEmpty()
    district: string;

    @Column()
    @IsNotEmpty()
    city: string;

    @Column()
    @MinLength(2)
    @MaxLength(2)
    @IsNotEmpty()
    state: string;

}